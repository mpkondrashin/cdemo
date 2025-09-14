from flask import Flask, request, render_template_string
import requests, json, time
import sys, subprocess

app = Flask(__name__)

def mcp_request(payload):
    """Один запрос к MCP серверу"""
    proc = subprocess.Popen(
        MCP_SERVER_CMD,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    out, err = proc.communicate(json.dumps(payload) + "\n", timeout=10)
    # Берём последнюю строку, которая является JSON
    for line in out.strip().splitlines()[::-1]:
        try:
            return json.loads(line)
        except:
            continue
    return {"error": "Invalid MCP response", "stdout": out, "stderr": err}

def list_tools():
    return mcp_request({
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/list"
    })

def call_tool(name, args):
    return mcp_request({
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {"name": name, "arguments": args}
    })

def ask_model(prompt, extra_context=""):
    """Вызов модели с добавлением данных из MCP"""
    payload = {
        "model": MODEL_ID,
        "prompt": f"{extra_context}\nUser: {prompt}\nAssistant:",
        "max_tokens": 200,
        "temperature": 0.7
    }
    r = requests.post(LOCALAI_BASE + "/v1/completions", json=payload, timeout=60)
    data = r.json()
    return data["choices"][0]["text"]


TEMPLATE = """
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>LocalAI Chat</title>
  <style>
    .spinner {
      border: 6px solid #f3f3f3;
      border-top: 6px solid #3498db;
      border-radius: 50%;
      width: 40px;
      height: 40px;
      animation: spin 1s linear infinite;
      margin: 10px auto;
    }
    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
  </style>
  <script>
    function showSpinner() {
      document.getElementById("response-area").innerHTML =
        '<div class="spinner"></div><div style="text-align:center">Loading...</div>';
    }
  </script>
</head>
<body>
  <h1>LocalAI Chat</h1>
  <form method="post" onsubmit="showSpinner()">
    <textarea name="prompt" rows="6" cols="70" placeholder="Type your message..."></textarea><br>
    <input type="submit" value="Ask">
  </form>

  <div id="response-area">
  {% if response %}
    <h2>Response:</h2>
    <div style="white-space:pre-wrap; border:1px solid #ccc; padding:8px;">{{ response }}</div>
  {% endif %}
  {% if error %}
    <h3 style="color:red;">{{ error }}</h3>
  {% endif %}
  </div>

  <p style="font-size:small">Model id: {{ model_id }}</p>
</body>
</html>
"""

LOCALAI_BASE = "http://127.0.0.1:8080"
# Try to read discovered model id from disk (written by userdata)
try:
    with open("/opt/local-ai/discovered_model_id","r") as f:
        MODEL_ID = f.read().strip()
except:
    MODEL_ID = ""

def get_model_id():
    global MODEL_ID
    if MODEL_ID:
        return MODEL_ID
    try:
        r = requests.get(LOCALAI_BASE + "/v1/models", timeout=5)
        data = r.json()
        MODEL_ID = data.get("data",[{}])[0].get("id","")
        # persist
        with open("/opt/local-ai/discovered_model_id","w") as f:
            f.write(MODEL_ID)
        return MODEL_ID
    except Exception as e:
        return MODEL_ID

@app.route("/", methods=["GET","POST"])
def chat():
    response_text = None
    error = None
    if request.method == "POST":
        user_prompt = request.form.get("prompt","")

        # 🔹 шаг 1: спросим у модели, нужно ли вызвать инструмент
        hint = ask_model(
            f"You have tools {list_tools()}. "
            f"If use asks for data from file system, pick tool and JSON parameters выбери. "
            f"Otherwise write 'none'.\n\Question: {user_prompt}\Reply:"
        )

        if "none" not in hint.lower():
            try:
                decision = json.loads(hint)  # {"tool":"list_files"} или {"tool":"read_file","args":{"filename":"x"}}
                tool = decision["tool"]
                args = decision.get("args", {})
                tool_result = call_tool(tool, args)
                # контекст для модели
                response_text = ask_model(user_prompt, extra_context=f"Результат инструмента {tool}: {tool_result}")
            except Exception as e:
                error = f"Ошибка разбора JSON из LLM: {hint} ({e})"
        else:
            response_text = ask_model(user_prompt)


@app.route("/prev", methods=["GET","POST"])
def chat_prev():
    error=None
    response_text=None
    model_id = get_model_id()
    if request.method == "POST":
        prompt = request.form.get("prompt","")

        if prompt.startswith("mcp:"):
            # MCP режим
            parts = prompt.split()
            cmd = parts[1] if len(parts)>1 else ""
            if cmd == "list_files":
                payload = {
                    "jsonrpc":"2.0","id":1,"method":"tools/call",
                    "params":{"name":"list_files","arguments":{}}
                }
                result = call_mcp(payload)
                response_text = json.dumps(result, indent=2, ensure_ascii=False)

            elif cmd == "read_file" and len(parts)>2:
                filename = parts[2]
                payload = {
                    "jsonrpc":"2.0","id":2,"method":"tools/call",
                    "params":{"name":"read_file","arguments":{"filename": filename}}
                }
                result = call_mcp(payload)
                response_text = json.dumps(result, indent=2, ensure_ascii=False)

            else:
                error = "Unknown MCP command"

        else:
            # Обычный запрос в LocalAI
            if not model_id:
                error = "Model not ready yet. Please wait and try again."
            else:
                payload = {
                    "model": model_id,
                    "prompt": prompt,
                    "max_tokens": 200,
                    "temperature": 0.7
                }
                try:
                    r = requests.post(LOCALAI_BASE + "/v1/completions", json=payload, timeout=120)
                    r.raise_for_status()
                    data = r.json()
                    if "choices" in data and len(data["choices"])>0:
                        response_text = data["choices"][0].get("text") or json.dumps(data)
                    else:
                        response_text = data.get("text") or json.dumps(data)
                except Exception as ex:
                    error = str(ex)

    return render_template_string(TEMPLATE, response=response_text, error=error, model_id=model_id)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=80)
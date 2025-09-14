import sys
import os
import json

FOLDER = "./data"  # папка, содержимое которой отдаём

def send_message(message):
    """Отправляем JSON сообщение в stdout"""
    data = json.dumps(message)
    sys.stdout.write(data + "\n")
    sys.stdout.flush()

def list_files():
    """Вернуть список файлов в папке"""
    try:
        return os.listdir(FOLDER)
    except Exception as e:
        return {"error": str(e)}

def read_file(filename):
    """Прочитать содержимое файла"""
    path = os.path.join(FOLDER, filename)
    if not os.path.isfile(path):
        return {"error": f"File {filename} not found"}
    try:
        with open(path, "r", encoding="utf-8") as f:
            return {"content": f.read()}
    except Exception as e:
        return {"error": str(e)}

def main():
    # handshake
    send_message({
        "jsonrpc": "2.0",
        "method": "server/ready",
        "params": {"message": "File MCP server started!"}
    })

    for line in sys.stdin:
        try:
            request = json.loads(line.strip())
            method = request.get("method")
            id_ = request.get("id")

            # список доступных инструментов
            if method == "tools/list":
                send_message({
                    "jsonrpc": "2.0",
                    "id": id_,
                    "result": {
                        "tools": [
                            {
                                "name": "list_files",
                                "description": "List files in the folder",
                                "input_schema": {
                                    "type": "object",
                                    "properties": {}
                                }
                            },
                            {
                                "name": "read_file",
                                "description": "Read contents of a file",
                                "input_schema": {
                                    "type": "object",
                                    "properties": {
                                        "filename": {"type": "string"}
                                    },
                                    "required": ["filename"]
                                }
                            }
                        ]
                    }
                })

            # вызов инструмента
            elif method == "tools/call":
                tool = request["params"]["name"]
                args = request["params"]["arguments"]

                if tool == "list_files":
                    send_message({
                        "jsonrpc": "2.0",
                        "id": id_,
                        "result": {"files": list_files()}
                    })

                elif tool == "read_file":
                    filename = args.get("filename")
                    send_message({
                        "jsonrpc": "2.0",
                        "id": id_,
                        "result": read_file(filename)
                    })

                else:
                    send_message({
                        "jsonrpc": "2.0",
                        "id": id_,
                        "error": {"code": -32601, "message": "Unknown tool"}
                    })

        except Exception as e:
            send_message({
                "jsonrpc": "2.0",
                "error": {"code": -32603, "message": str(e)}
            })

if __name__ == "__main__":
    main()

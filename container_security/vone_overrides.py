#!/usr/bin/env python3
import requests
import os
import yaml


class VOne:
    def __init__(self, token, region = ''):
        self.url_base = 'https://api.xdr.trendmicro.com'
        if region:
            self.url_base = f'https://api.{region}.xdr.trendmicro.com'
        self.token = token
    
    def request(self, method, url_path, params = None, headers = None):
        headers = headers or {}
        headers['Authorization'] = 'Bearer ' + self.token
        r = requests.request(method, url_path, params=params, headers=headers)
        if r.status_code >= 300:
            raise Exception(f'Error: {r.status_code}: {r.text}')
        
        if 'application/json' not in r.headers.get('Content-Type', ''):
            raise Exception(f'Unexpected content type: {r.headers.get("Content-Type")}')
        
        if not len(r.content):
            return None
            
        return r.json()

    def get_k8s_clusters(self, filter = ''):
        url_path = self.url_base +'/v3.0/containerSecurity/kubernetesClusters'
        query_params = {
            'top': '100',
            'orderBy': 'createdDateTime desc'
        }
        if filter:
            query_params['TMV1-Filter'] = filter
        while True:
            r = self.request('GET', url_path, params=query_params)
            data = r.json()
            for cluster in data['items']:
                yield cluster
                
            if 'nextLink' in data and data['nextLink']:
                url_path = data['nextLink']
            else:
                break
    
    def register_k8s_cluster(self, cluster_name, description, group_id):
        url_path = self.url_base +'/v3.0/containerSecurity/kubernetesClusters'
        body = {
            'name': cluster_name,
            'groupId': group_id,
            'description': description,
        }   
        r = self.request('POST', url_path, json=body)
        data = r.json()
        return data['apiKey'], data['endpointUrl']

    def iterate_cluster_groups(self, filter = ''):
        url_path = self.url_base +'/v3.0/containerSecurity/kubernetesClusterGroups'
        query_params = {
            'top': '100',
        }
        if filter:
            query_params['TMV1-Filter'] = filter
        r = self.request('GET', url_path, params=query_params)
        data = r.json()
        for cluster_group in data['items']:
            yield cluster_group

    def get_cluster_group(self, name):
        filter = f"name eq '{name}'"
        for cluster_group in self.iterate_cluster_groups(filter):
            return cluster_group
    
    def delete_k8s_cluster(self, id):
        url_path = self.url_base +f'/v3.0/containerSecurity/kubernetesClusters/{id}'
        self.request('DELETE', url_path)


cluster_name = 'demo_eks_cluster'

def write_overrides(api_key, endpoint_url):
    overrides = {
        'cloudOne': {
            'apiKey': api_key,
            'endpoint': endpoint_url,
            'exclusion': {
                'namespaces': ['kube-system']
            },
            'runtimeSecurity': {
                'enabled': True
            },
            'vulnerabilityScanning': {
                'enabled': True
            },
            'malwareScanning': {
                'enabled': True
            },
            'secretScanning': {
                'enabled': True
            },
            'inventoryCollection': {
                'enabled': True
            }
        }
    }
    with open('overrides.yaml', 'w') as f:
        f.write(yaml.dump(overrides))


def cleanup(token, region = ''):
    vone = VOne(token, region)
    for cluster in vone.get_k8s_clusters(f"name eq '{cluster_name}'"):
        print(f"[CDEMO] Deleting {cluster['name']}")
        vone.delete_k8s_cluster(cluster['id'])


def setup(token, region = ''):
    vone = VOne(token, region)
    print("[CDEMO] Getting existing clusters")
    for cluster in vone.get_k8s_clusters():
        print(f"[CDEMO] {cluster['name']}: {cluster['id']}")

    print("[CDEMO] Getting cluster groups")
    for group in vone.iterate_cluster_groups():
        print(f"[CDEMO] {group['name']}: {group['id']}")
    aws_eks_group = vone.get_cluster_group('Amazon EKS')
    if not aws_eks_group:
        raise Exception('Amazon EKS cluster group not found')
    group_id = aws_eks_group['id']
    api_key, endpoint_url = vone.register_k8s_cluster('demo_eks_cluster', 'Demo', group_id)
    print(f"[CDEMO] API Key: {api_key}")
    print(f"[CDEMO] Endpoint URL: {endpoint_url}")
    print("[CDEMO] Writing overrides.yaml")
    write_overrides(api_key, endpoint_url)
    print("[CDEMO] Done")


if __name__ == '__main__':
    api_token = os.getenv('API_TOKEN')
    if not api_token:
        raise Exception('API_TOKEN environment variable is not set')
    if os.argv[1] == 'cleanup':
        cleanup(api_token, os.argv[2])
    else:
        setup(api_token, os.argv[2])

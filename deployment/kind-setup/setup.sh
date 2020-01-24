kind delete cluster
kind create cluster
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-rc2/aio/deploy/recommended.yaml
kubectl apply -f dash-admin.yaml
kubectl apply -f dash-admin-bind.yaml
kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
kubectl label namespace default lightscreen=enabled
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}')

echo INFO run kubectl proxy
echo INFO go to 'http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login'

# for webhook:
# k apply -f deployment/webhook-registration.yaml
# don't forget to set up local ngrok and update the url
# test: kd deployment nginx && kubectl run nginx --image=nginx --replicas=1

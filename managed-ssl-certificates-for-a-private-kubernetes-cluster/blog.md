# Managed SSL Certs for a Private Kubernetes Cluster with CloudFlare, Cert Manager, and Let's Encrypt

Secure Socket Layer (SSL) certifications play a crucial role in your on-premise or cloud Kubernetes security. These certification:

1. Enable the ability to have encrypted traffic via the Transport Socket Layer (TSL) protocol between your Kubernetes cluster and a client device.
1. Ensure the integrity of the data being sent between the client and your Kubernetes cluster. 
1. Provide the client the ability to verify the identity of the service they are trying to communicate with.  

Even with all of these benefits a lot of companies still choose to not set up SSL certifications for their privately networked Kubernetes clusters. Managing SSL certificates from a trusted Certificate Authority (CA) traditionally is a cumbersome task, but with modern tools like [cert-manager](https://cert-manager.io/) and [Nginx Proxy Manager](https://nginxproxymanager.com/) it's a breeze!

In this tutorial, I am going to walk through how to set up cert-manager in a Minikube Kubernetes cluster to create and renew certifications automatically. You will also set up a `hello-world` deployment and service to test that we can receive HTTPS traffic via a Kubernetes Ingress. Grab some coffee and let's get started! ☕️

## Software Setup
For this tutorial, you will need to be running an amd64 or arm64 computer with either Linux or macOS as your operating system. Additionally, you will need to have [Minikube](https://minikube.sigs.k8s.io/docs/start/), [Docker](https://www.docker.com/products/docker-desktop/), and [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) installed. Once these three programs are installed, run the following command to create a Minikube cluster:

``` bash
minikube start --driver=docker --addons=ingress
```

This command starts minikube using docker as a driver. It also installs the needed services, RBAC roles, and deployments for an nginx Ingress Controller. You can validate that the Minikube node is up and running by using: 

``` bash
kubectl get nodes
```

If the node is online, you should `Ready` as the `Status` column in the following output:

```
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   64s   v1.27.4
```

## Domain and DNS Configuration

You will need a Cloudflare account and a domain where the nameserver (NS) records are set to be Cloudflare's. If you purchased your domain through Cloudflare, this step is irrelevant! To do this for a domain purchased with GoDaddy, go to your domain's DNS settings, click on the `Nameservers` tab and add the following entries:

```
damien.ns.cloudflare.com
davina.ns.cloudflare.com
```

It should look like: 

![GoDaddy Nameserver Setup](https://raw.githubusercontent.com/kevinlutzer/blogs/master/managed-ssl-certificates-for-a-private-kubernetes-cluster/assets/godaddy_nameservers.png "GoDaddy Nameserver Setup")

The steps to set the NS records are identical for other domain providers like Google. You will have to wait up to 48 hours for your domain's NS records to propagate to other nameservers. Additionally, you could set up your computer's DNS server to be Cloudflare's `1.1.1.1` server so you don't have to wait. To do this on macOS, navigate to your network details and then click on the `DNS` tab. Click on the `+` icon under `DNS Servers` and add `1.1.1.1`. The following image shows what this should look like:

![macOS DNS Setup](https://raw.githubusercontent.com/kevinlutzer/blogs/master/managed-ssl-certificates-for-a-private-kubernetes-cluster/assets/macos_dns_setup.png "macOS DNS Setup")

Set up a new A record for your domain in Cloudflare to point `hello-world.<YOUR DOMAIN>` to `127.0.0.1`. This is what the setup should look like in Cloudflare's UI: 

![Cloudflare Domain Setup](https://raw.githubusercontent.com/kevinlutzer/blogs/master/managed-ssl-certificates-for-a-private-kubernetes-cluster/assets/cloudflare_domain_setup.png "Cloudflare Domain Setup")

Pointing a DNS record to a virtual IP like `127.0.0.1` or a private IP like `192.168.1.2` is unusual, but it gives us the following benefits: 

1. Instead of having to modify your client device's host mapping in `/etc/hosts` or setting up a private DNS server, you can use Cloudflare's public DNS server.
1. We can get valid SSL certifications for a trusted Certificate Authority via a DNS-01 challenge type. More on this later.

## Cert Manager

cert-manager will be used in this tutorial to manage your SSL certifications. In your Kubernetes cluster, cert-manager is responsible for creating SSL certificates as well as renewing them before they expire. When installed, cert-manager adds two new Kubernetes object types that are used for creating and renewing certifications; the Issuer and Certificate.

Issuer objects are used to describe who the trusted CA you are using is, as well as the method used to obtain the certificate. For this tutorial, you will be using an [Automatic Certificate Management Environment](https://datatracker.ietf.org/doc/html/rfc8555) (ACME) issuer type with a DNS-01 solver. Let's Encrypt provides an ACME server that we can use to automate the domain ownership validation process to get our SSL certificates. There will be more on this later. 

The Certificate Object provides cert-manager information about when to renew the certificate, the expiration date, what encryption methods the certificate should use for the public/private key, what domain the certificate is for, and the owner email address.

To issue a certificate from Let's Encrypt, cert-manager gets a challenge value from Let's Encrypt's ACME servers. This value is to be put as a TXT record on your domain in Cloudflare. Once cert-manager creates this TXT record, it tells Let's Encrypt to validate the record's contents. Let's Encrypt compares the value of TXT record to that given to cert-manager. If they are the same, Let's Encrypt will issue a certificate to cert-manager. cert-manager will then create a Kubernetes secret from that certificate. Here is a diagram of how this all works, taken from an nginx [blog](https://www.nginx.com/blog/automating-certificate-management-in-a-kubernetes-environment/).

![DNS 01 Process](https://raw.githubusercontent.com/kevinlutzer/blogs/master/managed-ssl-certificates-for-a-private-kubernetes-cluster/assets/dns01_process.svg "DNS 01 Process")

For private and virtual IPs the DNS-01 challenge type is perfect as you don't need to have incoming traffic from the internet to your Kubernetes cluster to complete the process.

Install cert-manager into your Kubernetes cluster by running:

``` bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml
```

If the install was successful after a few minutes cert-manager should be running in the `cert-manager` namespace with three pods. Validate that the pods are up by running:

``` bash 
kubectl get pods -n cert-manager
```

You should see a similar output too: 

```
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-75d57c8d4b-dxkc6              1/1     Running   0          11s
cert-manager-cainjector-69d6f4d488-9bdst   1/1     Running   0          11s
cert-manager-webhook-869b6c65c4-9t7sb      1/1     Running   0          11s
```

### Creating the Cloudflare API Token

Now you need to create a secret containing a Cloudflare API token. This token will be used to create and delete the TXT record on your domain needed during the DNS-01 challenge process. To do this, navigate to https://dash.cloudflare.com/profile/api-tokens. From there, click on `Create Token` and then on the following page click on the `Use Template` button in the same row as `Edit zone DNS`. This will allow you to configure what zone resources the token will have access to. Under Zone Resources choose `Include All Zones`. You can additionally add `Client IP Address Filtering` as well. This is **optional**, but I recommend this as it helps add some protection for your token if it is leaked. This IP will be the public IP of your computer. Click on continue to summary and the `Create Token`. Take the token value you get from that page and create a Kubernetes secret from it using the following command by replacing `<TOKEN>` with your token: 

``` bash
kubectl create secret generic cloudflare-api-key-secret --from-literal=api-key=<TOKEN>
```

### Issuer
Next, you need to create the `issuer.yml` file. Add the following contents to that file but replace the email values with your Cloudflare account's email address.

``` yaml
# issuer.yml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
spec:
  acme:
    email: <YOUR EMAIL>
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: issuer-key
    solvers:
    - dns01:
        cloudflare:
          email: <YOUR EMAIL>
          apiTokenSecretRef:
            name: cloudflare-api-key-secret
            key: api-key
```

Create the Issuer object by running:

`kubectl apply -f issuer.yml`

### Certificate

Create a new file with the following contents, but with substituting `<YOUR DOMAIN>` for your domain. 

``` yaml 
# certificate.yml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: hello-world-ca-tsl
spec:
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  subject:
    organizations:
      - testorganization
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048
  dnsNames:
    - hello-world.<YOUR DOMAIN>
  secretName: hello-world-ca-tsl
  issuerRef:
    name: ca-issuer
    kind: Issuer
    group: cert-manager.io

```

Create the certificate object:

``` bash
kubectl apply -f certificate.yml
```

Immediately after the Certificate object is created, cert-manager will start the process of creating the certificate. Verify that the certificate was successfully issued by running:

``` bash 
kubectl get certificate
```

This **will take a few minutes**, but once the certificate has been successfully issued, you will see the `READY` column read `True` like in the following output:

```
NAME                 READY   SECRET               AGE
hello-world-ca-tsl   True    hello-world-ca-tsl   100s
```

If you would like to view the contents of this newly created certificate, run the following commands:

``` bash
kubectl get secret foo-ca-tsl -o jsonpath='{.data.*}' | base64 -d >> cert.crt
openssl x509 -in cert.crt -text -noout 
```

## Hello World Service Setup

Now that we have a valid SSL certificate let's set up a demo Deployment, Service and Ingress to test that you can make a HTTPS request with that certificate's contents. Run the following command to create the `hello-world` deployment.

``` bash 
kubectl create deployment hello-world --image=klutzer/hello-world:latest
```

The image specified sets up a web server that responds to requests on port 80 with `Hello, World!😀`. Now expose this deployment through a service:

``` bash 
kubectl expose deployment hello-world --type=NodePort --port=80
```

Verify the deployment and service are working by running: 

``` bash
kubectl get pods,service
```

You should see an output similar to the following:

```
NAME                             READY   STATUS    RESTARTS   AGE
pod/hello-world-bbbbdbfc-g8jrq   1/1     Running   0          14s

NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
service/hello-world   NodePort    10.110.13.203   <none>        80:30639/TCP   10s
service/kubernetes    ClusterIP   10.96.0.1       <none>        443/TCP        48m
```

Create an `ingress.yml` file with the following contents, replacing `<YOUR DOMAIN>` with your domain:

``` yaml
# ingress.yml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-world-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  tls:
    - hosts:
      - hello-world.<YOUR DOMAIN>
      secretName: hello-world-ca-tsl
  rules:
    - host: hello-world.<YOUR DOMAIN>
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-world
                port:
                  number: 80
```

This ingress will be used to terminate HTTPS requests on `hello-world.<YOUR DOMAIN>` and pass on HTTP traffic to your service. Terminating HTTPS requests with an Ingress is a very common way to provide HTTPS access to your services. This gives you the benefit of not having to build SSL support for your web servers. If you want to learn more about Kubernetes Ingresses check out the documentation [here](https://kubernetes.io/docs/concepts/services-networking/ingress/). Run the following command to create the ingress file.

``` bash 
kubectl apply -f ingress.yml
```

Since Minikube runs in a docker container, you need to be able to port forward into it to get access to the `hello-world` service. To do this, run the following command in a new terminal. 

``` bash
sudo minikube tunnel 
```

**Note**: This command will need to run as long as you want to access the services in your Minikube cluster. To test that you can call the `hello-world` service, run:

``` bash 
curl https://hello-world.<YOUR DOMAIN>
```

You should see the following response:

```
Hello, World!😀
```

Thats it! You have set up cert-manager and created an SSL certificate for a simple `hello-world` application. 

## Conclusion

cert-manager is a powerful tool for creating and renewing certificates. With Let's Encrypt and the DNS-01 challenge type we can create SSL certificates for Kubernetes services running on private networks. This allows you to maintain the security of the traffic between clients and your Kubernetes services with minimal effort. 

## Teardown

To stop and delete your Minikube cluster run the following commands: 

``` bash
minikube stop
minikube delete --all
```

Make sure you delete the Cloudflare API Token as well!

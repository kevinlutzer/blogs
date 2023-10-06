# Managed TSL Certs for An On Premises Kubernetes Cluster with CloudFlaire, Cert Manager, and Let's Encrypt

Secure Socket Layer (SSL) certifications play a curcial role in your on premise or cloud security. These certification:

1. Enable the ability to have encrypted traffic using HTTPs between your Kubernetes cluster and a client device.
1. Ensure the integrity of the data being sent between the client and your Kubernetes cluster. 
1. Provide the client the ability to verify the identity of the service they are trying to communicate with.  

Even with all of these benefits a lot of company still choose to not setup TSL certifications for their private Kubernetes clusters or on premises services. Traditionally managing certifications manually can be a bothersome task, but with modern automation tools like [cert-manager](https://cert-manager.io/) and [Nginx Proxy Manager](https://nginxproxymanager.com/) it is a breeze! In this tutorial I am going talk about what is in a TSL certification, what the types of Certifications are, how to setup cert-manager in your Kubernetes cluster to create and update certifications automagically, and how to use the certification to provide HTTPS terminmation for an nginx ingress exposing a hello world service. Grab some coffee and lets's get started! ☕️


`https://dash.cloudflare.com/profile/api-tokens`


``` bash 
curl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
     -H "Authorization: Bearer <YOUR TOKEN>" \
     -H "Content-Type:application/json"
```

## About Certificates

Simply, a TSL certificate provides information about: 

1. The identity of the certificate creator creator.
1. Information about the domain or IP of the server the certificate is for.
1. The public and private keys that the server will use to decrypt traffic encrypted over TSL.
1. Expiration information

Every TSL certificate needs to have its information signed/created by a Certificate Authorities's (CA) certificate. The CA is responsible for validating that the domain or IP the certificate is being created for is correct, as well as the encryption keys used are valid. How do we know if a CA actually validated these details for a certificate our computer gets when trying to access a server online? The answer to that question is based on if the certificate is created by a trusted CA or not. Trusted Certificate Authorities include orginazations like Google, Let's Encrypt and Cloudflare. When creating a TSL certificate with one of them, you have to prove you control the domain the certificate is for. If you can prove you control the domain, the CA will issue you a TSL certificate!

<MORE ABOUT ISSUING EHRE>

## Software Setup

For this tutorial you will need to be running a amd64 or arm64 (Apple M1, M2) based computer with either Linux or macOS as your operating system. Additionally you will need to have [Minikube](https://minikube.sigs.k8s.io/docs/start/), [Docker](https://www.docker.com/products/docker-desktop/), and [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) installed. Once these three programs are installed, run the following command to create a minikube cluster:

``` bash
minikube start --driver=docker --addons=ingress
```

You can validate that the Minikube node is up and running by using: 

``` bash
kubectl get pods
```

If the node is ready, you should see the following output:

```
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   12m   v1.27.4
```

## Domain and DNS Configuration

You will need a Cloudflare account and a domain where the nameservers are set to be Cloudflare's. To do this for a domain purchased with Godaddy, go to your domain's DNS settings, click on the `Nameservers` tab and add the following entries

```
damien.ns.cloudflare.com
davina.ns.cloudflare.com
```

It should look like: 

<INSERT IMAGE>

If you just set the NS records for your domain, you will either have to wait about 24 hours before you continue, or set your computer's DNS server to be Cloudflare's `1.1.1.1` server. Setup a new A record for your domain in Cloudflare to point `hello-world.<YOUR DOMAIN>` to `127.0.0.1`. This will enable all local HTTPS traffic on your computer to `hello-world`.<YOUR DOMAIN> to go to your computer's virtual IP of 127.0.0.1. 

**Note**: This will not cause an issue with certification issuing, as the issuing challenge you will use is DNS-01. Unlike traditional HTTP challenges, DNS-01 does not require Let's Encrypt to hit `127.0.0.1`. 

## Cert Manager

You will be using cert-manager to mange your TSL certifications. Cert manager creates two new Kubernetes object types that are used for creating and updating certification; the Issuer and Certificate. Issuer objects are used to describe who the CA is, as well as the method used to obtain the certificate. For this tutorial you will be using a DNS-01 challenge to secure a TSL certificate from Let's Encrypt.

To issue a certificate from Let's Encrypt, cert-manager gets a challenge value from Let's encrypt's servers to put as a TXT DNS record on your domain in Cloudflare. Once cert-manager creates this TXT record, it tells cloudflare to validate the record's contents. Let's encrypt validates the TXT record and if everything is correct, will issue a certificate to cert-manager. cert-manager will then create a Kubernetes secret with the contents being the certificate. Here is a diagram of how this all works taken from an nginx [blog](https://www.nginx.com/blog/automating-certificate-management-in-a-kubernetes-environment/)

<INSERT IMAGE>

Install cert-manager into your Kubernetes cluster by running:

``` bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml
```

Validate that the pods are up by running:

``` bash 
kubectl get pods
```

If the install was successful you should see three pods with the status "Running".

```
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-75d57c8d4b-trnsr              1/1     Running   0          8m36s
cert-manager-cainjector-69d6f4d488-fcpfd   1/1     Running   0          8m36s
cert-manager-webhook-869b6c65c4-fnjmp      1/1     Running   0          8m36s
```

Now you need to create a secret containing a Cloudflare API token that will be used to created the DNS TXT record on your domain. To do this, navigate to: 

https://dash.cloudflare.com/profile/api-tokens

From there, click on `Create Token` and then on the following page click on the `Use Template` button in the same row as `Edit zone DNS`. This will allow you to configure what zone resources the token will have access too. Under Zone Resources choose `Include All Zones`. You can additionally add `Client IP Address Filtering` as well. I recommend this as it helps add some protection for your token if it is accidentally leaked. You can get your computer's public IP by going to https://whatismyipaddress.com/. Click on continue to summary and the `Create Token`.

Take the value you get from that page and create a Kubernetes secret using the following command: 

``` bash
kubectl create secret generic cloudflare-api-key-secret --from-literal=api-key=<TOKEN>
```
### Issuer

Next you need to create the `issuer.yml` file. Add the following contents to that file but replace the email values with your Cloudflare account's email address.

``` yaml
# issuer.yml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
spec:
  acme:
    email: <YOUR CLOUDFLARE EMAIL>
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: issuer-key
    solvers:
    - dns01:
        cloudflare:
          email: <YOUR CLOUDFLARE EMAIL>
          apiTokenSecretRef:
            name: cloudflare-api-key-secret
            key: api-key
```

The last step needed to get a certificate from Let's encrypt is to create a cert-manager Certificate object. Create a new file with the following contents with replacing the dnsNames to be your own. 

### Certificate

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

Create the certificate object

``` bash
kubectl apply -f certificate.yml
```

Verify that the certificate was successfuly issued by running:

``` bash 
kubectl get certificate
```

This will take a few minutes, but once the certificate has been successfully issued and the Kubernetes secret containing it's value has been created, you will see the `READY` column read `True`.  
 
```
NAME         READY   SECRET                   AGE
hello-world-ca-tls   True    hello-world-ca-tls   2m14s
```

If you would like to view the contents of this newely created certificate, run the following commands:

``` bash
kubectl get secret foo-ca-tsl -o jsonpath='{.data.*}' | base64 -d >> cert.crt
openssl x509 -in cert.crt -text -noout 
```

## Hello World Setup
Now that we have a valid TSL certificate lets setup a demo Deployment, Service and Ingress to test that you can make a HTTPS request with that certificate's contents. Run the following command to create a deployment where the container is running a webserver that responds to requests on port 80 with `Hello, World!😀`.

``` bash 
kubectl create deployment hello-world --image=klutzer/hello-world:latest
```

Now expose this deployment through a service:

``` bash 
kubectl expose deployment hello-world --type=NodePort --port=80
```

Verify the deployment and service are working by running: 

``` bash
kubectl get pods,service
```

Create an `ingress.yml` file with the following contents. This ingress will be used to terminate HTTPs requests on hello-world.<YOUR DOMAIN>.ca and pass the unecrypted HTTP traffic to the hello world Kubernetes Pods. 

``` yaml
# ingress.yml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-world-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
    - host: hello-world.<YOUR DOMAIN> # example: hello-world.kevinlutzer.ca
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

Run the following command to create the ingress file 

``` bash 
kubectl apply -f ingress.yml
```

Now in one terminal, run: 

``` bash
minikube tunnel 
```

This will create allow you to port foward and access the hello world Ingress in the Minikube cluster. To test the that you can call the hello world service run:

``` bash 
curl https://hello-world.<YOUR DOMAIN>
```

You should see the following response

```
Hello, World!😀
```

## Teardown

To clean up your minikube cluster
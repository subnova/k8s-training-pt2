# Kubernetes training part 2

## Start minikube

```sh
minikube start
```

## Build our test app

```sh
cd app
minikube image build -t configtest .
cd ..
```

## Deploy our test app

```sh
$ kubectl create ns demo
$ kubens demo
$ kubectl apply -f 00-default.yaml
$ kubectl apply -f service.yaml
$ export URL=http://$(minikube ip):$(kubectl get svc/configtest --output='jsonpath={.spec.ports[0].nodePort}')/hello
$ curl $URL
Hello world
```

## Set a command line argument

```sh
$ kubectl apply -f 01-args.yaml
$ curl $URL
Hello Testers
```

## Set an environment variable

```sh
$ kubectl apply -f 02-env.yaml
$ curl $URL
Hello Developers
```

## Set a command line argument from an environment variable

```sh
$ kubectl apply -f 03-args-from-env.yaml
$ curl $URL
Hello Developers!
```

## Set an environment variable using the Downward API

```sh
$ kubectl apply -f 04-env-from-downward.yaml
$ curl $URL
Hello configtest-xxxxxxxxx-yyyyy
```

## Set an environment variable from other environment variables

```sh
$ kubectl apply -f 05-env-from-env.yaml
$ curl $URL
Hello Developers and Testers from configtest-xxxxxxxxx-yyyyy
```

## Set environment variables from a configmap

```sh
$ kubectl apply -f 06-env-from-configmap.yaml
$ kubectl get po
NAME                        READY   STATUS                      RESTARTS    AGE
configtest-xxxxxxxxx-yyyyy  0/1     CreateContainerConfigError  0           1m
configtest-ppppppppp-qqqqq  1/1     Running                     0           5m
$ kubectl apply -f configmap.yaml
$ kubectl get po
NAME                        READY   STATUS                      RESTARTS    AGE
configtest-xxxxxxxxx-yyyyy  1/1     Running                     0           5m
$ curl $URL
Hello UX
```

## Set environment variables from a configmap reference

```sh
$ kubectl apply -f 07-env-from-configmap-element.yaml
$ curl $URL
Hello Designers
$ kubectl edit cm/configtest # Replace greetee value with "Product Folk"
$ curl $URL
Hello Designers
$ kubectl apply -f 07-env-from-configmap-element.yaml
$ curl $URL
Hello Designers
$ kubectl delete $(kubectl get po -o name)
$ curl $URL
Hello Product Folk
```

## Set environment variables from a secret

```sh
$ kubectl apply -f secret.yaml
$ kubectl apply -f 08-env-from-secret.yaml
$ curl $URL
Hello Grand Master
```

## Set volume from a config map

```sh
$ curl -H 'Accept-Language: fr' $URL
Bonjour world
$ curl -H 'Accept-Language: de' $URL
Hello world
$ kubectl apply -f 09-volume-from-configmap.yaml
$ kubectl get po
NAME                        READY   STATUS                      RESTARTS    AGE
configtest-xxxxxxxxx-yyyyy  0/1     ContainerCreating           0           1m
configtest-ppppppppp-qqqqq  1/1     Running                     0           5m
$ kubectl describe $(kubectl get po --field-selector status.phase=Pending -o name)
$ kubectl apply -f langmap.yaml
$ kubectl get po
NAME                        READY   STATUS                      RESTARTS    AGE
configtest-xxxxxxxxx-yyyyy  0/1     Running                     0           1m
$ curl -H 'Accept-Language: de' $URL
Gutten tag world
```

## Set volume from a secret

```sh
$ kubectl apply -f 10-volume-from-secret.yaml
$ curl $URL
$ curl -H 'Accept-Language: il' $URL
Hello Illuminatus Grand Master 👁⃤
```

## Templating with Helm

```sh
$ cd configtest
$ helm upgrade --install --create-namespace --namespace helmdemo configtest .
$ kubens helmdemo
$ helm list
NAME            NAMESPACE       REVISION        UPDATED                                STATUS          CHART                   APP VERSION
configtest      helmdemo        1               2021-10-08 10:03:38.2957355 +0100 BST  deployed        configtest-0.1.1        1.0.0   
$ kubens helmdemo
$ kubectl get po
$ export URL=http://$(minikube ip):$(kubectl get svc/configtest --output='jsonpath={.spec.ports[0].nodePort}')/hello
$ curl $URL
Hello Pags
$ helm upgrade --install --create-namespace --namespace helmdemo --set greetingName=Shaheen configtest .
$ curl $URL
Hello Shaheen
$ echo "greetingName: Renuka" > override.yaml
$ helm upgrade --install --create-namespace --namespace helmdemo --values override.yaml configtest .
$ curl $URL
Hello Renuka
$ cd ..
```

## Kustomization

```sh
$ cd kustomize
$ kubectl apply -k base
$ kubens kustomconfig
$ kubectl get po
$ export URL=http://$(minikube ip):$(kubectl get svc/configtest --output='jsonpath={.spec.ports[0].nodePort}')/hello
$ curl $URL
Hello world
$ GREETING_NAME=Stuart kubectl apply -k base
$ curl $URL
Hello Stuart
$ kubectl describe po
$ kubectl get cm
NAME                    DATA   AGE
configtest-xxxxxxxxxx   1      2m
configtest-yyyyyyyyyy   1      5m
$ kubectl rollout undo deploy/configtest
$ curl $URL
Hello world
$ curl -H 'Accept-Language: de' $URL
Hello world
$ kubectl apply -k overlays/prod
$ kubens prodconfig
$ kubectl get po
$ export URL=http://$(minikube ip):$(kubectl get svc/configtest-prod --output='jsonpath={.spec.ports[0].nodePort}')/hello
$ curl -H 'Accept-Language: de' $URL
Gutten tag world
$ GREETING_NAME=Shrihari kubectl apply -k overlays/prod
$ curl -H 'Accept-Language: de' $URL
Gutten tag Shrihari
$ cd ..
```

## Skaffolding

```sh
$ skaffold dev --port-forward
$ curl http://localhost:8090/hello
Hello world
$ git apply main.go.patch
$ curl http://localhost:8090/hello
Hello world!
$ ^C
$ skaffold run --profile prod
$ kubens prodconfig
$ kubectl get po
```
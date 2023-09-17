# Grupo27-Laboratorio-1

# **Integrantes**

---

- Nicolas Barrera 201951552-7

- Daniel  Sepúlveda 201873065-3

- Javier Maturana 201873604-k

# **Instrucciones**

---

**Primero es necesario para que todas las VM esten en la 
ruta correcta es necesrio ejecutar el siguiente comando en
todas las VM:**


```
cd Grupo27-Laboratorio-1
```



1. El rabbit se inicia en la dist106 usando el comando:

```
make docker-rabbit
```

2. Despues es necesario ejecutar en todas las VM el comando: 

```
make build
```

3. Luego en todas las VM se debe ejecutar: 

```
make docker-regional
```

4. Despues en la dist105 se debe ejecutar:

```
make docker-central
```

5. Por ultimo para revisar los los log de las VM regional y/o central se debe 
ejectuar en la dist105:

```
docker logs-central
```
6. Una vez terminadas todas las iteraciones:
```
Cerrar con control-c
```


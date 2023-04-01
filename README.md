# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un ejemplo de cliente-servidor el cual corre en containers con la ayuda de [docker-compose](https://docs.docker.com/compose/). El mismo es un ejemplo práctico brindado por la cátedra para que los alumnos tengan un esqueleto básico de cómo armar un proyecto de cero en donde todas las dependencias del mismo se encuentren encapsuladas en containers. El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers.

Por otro lado, se presenta una guía de ejercicios que los alumnos deberán resolver teniendo en cuenta las consideraciones generales descriptas al pie de este archivo.

## Resolución 

### Scripts disponibles

* ``./build.sh [server|client]`` : Crea las imagenes utilizando la configuración definida en cada respectivo Dockerfile.
* ``./run.sh n`` : Levanta docker compose con los recursos configurados según 'tp0-doccgen.py' utilizando "n" clientes.
* ``./run.sh logs`` : Abre los logs de los servicios.
* ``./run.sh test`` : Con el servidor levantado, crea un container temporal dentro de la misma red para probar el servidor.
  * Dentro de la instancia de prueba ejecutar ``./test-server.sh``.
* ``./stop.sh [-k]`` : Detiene los recursos creados en docker compose. Al utilizar "-k" los elimina.

### Protocolo de comunicación utilizado

El protocolo de comunicación implementado se basa en el pasaje de mensajes sobre el protocolo TCP de la capa de transporte, utilizando paquetes inspirados en el formato TLV (Type-Lenght-Value) para soportar operaciones con mensajes variables y así hacer un uso justo y necesario de la red.

El protocolo implementado se divide en 3 capas:

3- Capa de aplicación  
2- Capa empaquetado ("tlv")     
1- Capa de comunicación         

#### Comunicación
Es la capa mas baja del protocolo implementado y se ubica inmediatamente sobre la capa de transporte (protocolo TCP).
Su función es proveer una primera abstracción de los sockets, implementando `send(buffer)` y `receive(n)` que permitan una escritura y lectura limpia, evitando fenómenos como short-read y short-write.

#### Empaquetado
Esta capa se ubica sobre la anterior y en ella se implementa el formato de mensajes TLV. Su función es organizar la lectura y escritura en dos etapas: una para el header y otra para el body. De esta manera se consiguen mensajes de longitud variable.

Se define un paquete de capa de aplicación llamado "Packet" que cuenta con los siguientes 3 campos:

- `OPCODE` - Es el código de operación (Type). Se asigna 1 byte para este campo y sus posibles valores son definidos por la capa superior que contendrá la lógica de negocio.
- `DATA_LENGHT` - Es el tamaño en bytes del body (Lenght). Se asignan 4 bytes para este campo.
- `DATA` - Cadena de bytes que será enviada o recibida por el socket (Value). Representa la información necesaria para ejecutar el OPCODE correspondiente.

Otra de las funciones importantes de esta capa es fragmentar el vector de bytes en varios segmentos del tamaño maximo configurado para ser enviados, y su contraparte, ensamblar los fragmentos recibidos en un único vector de bytes para luego ser enviado a la capa superior.

#### Aplicación 
Esta capa define cuales son los OPCODES que serán válidos y provee interfaces para que el cliente pueda hacer uso de estas operaciones abstrayendose del mecanismo de comunicación y una interfaz para que el servidor pueda procesar las operaciones.

Códigos de operacion definidos:

| Codigo  | Función | Data |
| ------------- | ------------- | ---- |
| -1 | Mensaje enviado/recibido corrupto | _empty_ |
| 0  | Cierre de conexión | _empty_ |
| 1  | Registrar una apuesta  | [ `id-cliente`, `nombre`, `apellido`, `documento`, `nacimiento`, `numero`] |
| 2 |  ACK del servidor  | _empty_
| 3  | Registrar un batch de _n_ apuestas  | [ n , "@0", [apuesta_0] , .. "@n", [apuesta_n]  ] |
| 4  | Comunicar un error | Mensaje de error |
| 5  | Agencia lista para sorteo | [ `id-client`] |
| 6  | Solicitar ganadores | [ `id-client` ] |
| 7  | Envío de ganadores | [ `dni1`, .. ,`dniN` ] |
| 8  | Servidor ocupado | _empty_ |


### Concurrencia
El servidor implementa procesamiento concurrente para poder atender las requests de los clientes y mecanismos de sincronización que permiten
el acceso controlado a recursos compartidos.

Para poder atender concurrentemente a los clientes se tiene un pool de `n` executors (parámetro configurable), objeto provisto por la
librería ``concurrent.futures`` que permite lanzar concurrentemente las tareas que se encargan de manejar las conexiones entrantes del servidor.


Por otro lado, el mainthread se encarga de continuar escuchando nuevas conexiones.

El mecanismo de sincronización utilizado principalmente para la escritura y lectura de archivos y el acceso a variables de control compartidas
es el `Lock()` que provee la librería `threading`. 

------

## Requisitos de ejecución

#### Ej5
Se debe configurar un archivo `.env` dentro de la carpeta ./client/ con la los valores de los campos a enviar.


#### Ej6
Se debe descomprimir el `dataset.zip` dentro de la carpeta ./.data/dataset/


---
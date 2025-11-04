# theNewWorldGame
¡Hola! Este es mi proyecto "The new world", desarrollado para mi materia de Programación Concurrente. Es un pequeño juego de simulación 2D hecho en Go con la librería EbitenEngine.

El objetivo principal no era solo hacer un juego, sino demostrar cómo aplicar los patrones de concurrencia de Go (goroutines, canales, Mutex) en una aplicación visual e interactiva que corre sin congelarse.

Características
Gestión de 5 Sobrevivientes: Dirige a 5 "survivors" (controlados por sprites) en el mapa.

Tareas Concurrentes: ¡Puedes asignarles tareas a los 5 al mismo tiempo! Cada uno trabaja en su propio hilo (goroutine) sin congelar el juego.

Recolección y Crafteo: Los survivors pueden recolectar Wood y Scrap (haciendo clic) o construir/mejorar el Refugio (tecla B) y Barreras (tecla V).

Gestión de Recursos: Un inventario centralizado (Store) que se actualiza de forma segura.

Visualización de Estados: Puedes ver en tiempo real si un survivor está IDLE, MOVING_TO_RESOURCE, GATHERING (trabajando), o MOVING_TO_BASE.

Mi Arquitectura (Separación de Responsabilidades)
Para cumplir el requisito, diseñé el proyecto en 3 capas claras:

Capa domain :

Aquí vive la lógica pura. Son solo structs como Task, Result, y lo más importante, el Store (Almacén).

Sincronización: Para protegerme de race conditions, metí el sync.Mutex dentro del Store. Así, el inventario se "protege solo", y cualquier otra capa que lo use (AddResource, ConsumeResources) ya lo está haciendo de forma segura.

Capa application:

Aquí vive toda la concurrencia, o sea la lógica.

Aquí definí la lógica de mis workers (survivor/survivor.go) y el GameService que los "contrata" y crea los canales (jobsChan, resultsChan).

Capa infra:

Este es el código de Ebiten (game.go)

NO hace trabajo pesado, solo dibuja los sprites, detecta clics, y se comunica con la "sala de máquinas" dándoles órdenes (enviando Task al jobsChan) y leyendo sus reportes (recibiendo Result del resultsChan).

1 .- Inicialización: El main.go primero crea el GameService  y luego se lo "inyecta" al struct Game de ebiten. 

2.- . Update() como productor: ya que este envía tareas. El hilo principal (Update) no envía la tarea a la goroutine inmediatamente, este gestiona los estados visuales, el visual visor , o sea solo cuando el sprite (mi personaje)  llega a su destino, el Update envía la tarea al jobsChan.

3. Update() como consumidor: El update revisa el results chan y si está vacío, el default se activa y el juego sigue corriendo, pero si hay un resultado, lo procesa.

4. Draw() : Finalmente el hilo draw donde se dibujado todo! E igualmente necesita leer el inventario, o sea la memoria compartida para dibujarlo, lo cual puede resultar con un problema, pero acá lo maneje para que se haga de forma segura llamando al método del Store, que usa el Mutex.

Patrón worker pool
El proyecto implementa el patrón Worker pool, en donde el pool, son las 5 goroutines SurvivorMainLoop lanzadas por NewGameService y la cola de tareas es el canal jobsChan.
En donde la lógica de game.go actúa como un productor que distribuye tareas (task) a la cola, las 5 goroutines del pool compiten como consumidores por así decirlo por tomar tareas de esa cola, entonces esto permite que el jugador asigne hasta 5 tareas (recolección, crafteo) y estas se ejecuten simultáneamente en el backend, mientras el uegue funcione sin problemas.


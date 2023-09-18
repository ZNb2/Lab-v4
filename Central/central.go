package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/MetalDanyboy/Lab1/protos"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Server_name(servidor string) (host string, nombre string, puerto string){
	
	//host = "localhost"
	if servidor == "America"{
		host = "dist105.inf.santiago.usm.cl"
		puerto = "50052"
		nombre = "America"
	}else if servidor == "Asia"{
		host = "dist106.inf.santiago.usm.cl"
		puerto = "50053"
		nombre = "Asia"
	}else if servidor == "Europa"{
		host = "dist107.inf.santiago.usm.cl"
		puerto = "50054"
		nombre = "Europa"
	}else if servidor == "Oceania"{
		host = "dist108.inf.santiago.usm.cl"
		puerto = "50055"
		nombre = "Oceania"
	}
	return
}

func ConexionGRPC2(keys int, servidor string, wg *sync.WaitGroup){
	
	host, nombre, puerto := Server_name (servidor)
	log.Println("Connecting to server "+nombre+": "+host+":"+puerto+". . .")
	conn, err := grpc.Dial(host+":"+puerto,grpc.WithTransportCredentials(insecure.NewCredentials()))	
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//fmt.Printf("Esperando\n")
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)
	for {
		log.Println("Sending message to server "+nombre+": "+strconv.Itoa(keys))
		_, err := c.SendKeys(context.Background(), &pb.NumberRequest{Number: int32(keys)})
		if err != nil {
			log.Println("Server "+nombre+" not responding: ")
			log.Println("Trying again in 10 seconds. . .")
			time.Sleep(10 * time.Second)
			continue
		}
		//log.Printf("Response from server "+nombre+": "+"%s", response.Body)
		break
	}
	defer wg.Done()
}

func ConexionGRPC(mensaje string, servidor string , wg *sync.WaitGroup){
	
	host, nombre, puerto := Server_name (servidor)
	log.Println("Connecting to server "+nombre+": "+host+":"+puerto+". . .")
	conn, err := grpc.Dial(host+":"+puerto,grpc.WithTransportCredentials(insecure.NewCredentials()))	
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//fmt.Printf("Esperando\n")
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)
	for {
		log.Println("Sending message to server "+nombre+": "+mensaje)
		_, err := c.SayHello(context.Background(), &pb.Message{Body: mensaje})
		if err != nil {
			log.Println("Server "+nombre+" not responding: ")
			log.Println("Trying again in 10 seconds. . .")
			time.Sleep(10 * time.Second)
			continue
		}
		//log.Printf("Response from server "+nombre+": "+"%s", response.Body)
		break
	}
	defer wg.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Starting Central. . .")

	directorioActual, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        return
    }
    content, err := os.ReadFile(directorioActual+"/Central/parametros_de_inicio.txt")
	if err != nil {
		log.Fatal(err)
	}
	lineas := strings.Split(string(content), "\n")
	rangoLlaves := strings.Split(lineas[0], "-")
	var min, max , iterations, contador int
	min, _= strconv.Atoi(rangoLlaves[0])
	max, _= strconv.Atoi(rangoLlaves[1])
	iterations, _= strconv.Atoi(lineas[1])

	

	


	//...CONEXION RABBITMQ...
	addr := "dist106.inf.santiago.usm.cl"
	//addr :="localhost"
	connection, err := amqp.Dial("amqp://guest:guest@"+addr+":5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		"testing", // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       //args
	)
	if err != nil {
		panic(err)
	}
	// ...
	
	var llaves int
	for {	
			contador++
			if iterations != -1 {
				if contador == iterations+1{
					break
				}
			}
			
			if iterations == -1 {
				fmt.Printf("\nGeneración %d/infinito\n", contador)
			}else{
				fmt.Printf("\nGeneración %d/%d\n", contador,iterations)
			}
			llaves= rand.Intn(max-min) + min

			log.Printf("Llaves disponibles: %d", llaves)
			
			
			var wg sync.WaitGroup
			wg.Add(1)
			go ConexionGRPC("LLaves Disponibles","America", &wg)
			wg.Add(1)
			go ConexionGRPC("LLaves Disponibles","Asia", &wg)
			wg.Add(1)
			go ConexionGRPC("LLaves Disponibles","Europa", &wg)
			wg.Add(1)
			go ConexionGRPC("LLaves Disponibles","Oceania", &wg)
			wg.Wait()
		
			
			
			//Mensaje Rabbit
			
			var num_cola int
			var wg2 sync.WaitGroup
			for msg := range msgs {
				     
					//fmt.Printf("Received Message: %s\n", msg.Body)
					subcadenas := strings.Split(string(msg.Body), "-")
					
					llaves_pedidas,_:=strconv.Atoi(subcadenas[1])
					if llaves_pedidas > llaves{
						llaves_pedidas=llaves
					}
					if llaves != 0{
						llaves-=llaves_pedidas
					}

					log.Printf("Mensaje asíncrono de servidor %s leído\n", subcadenas[0])
					wg2.Add(1)
					go ConexionGRPC2(llaves_pedidas,subcadenas[0], &wg2)
					num_cola++
					log.Printf("Se inscribieron %d cupos de servidor %s\n", llaves_pedidas, subcadenas[0])
					
					if num_cola == 2{
						break
					}	
			}
			wg2.Wait()
			//time.Sleep(1 * time.Second)
		}
	defer log.Println("Closing Central. . .")
}







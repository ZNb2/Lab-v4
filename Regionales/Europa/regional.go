package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/MetalDanyboy/Lab1/protos"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)
var server_name string
var cant_registrados int
var cant_llaves_pedidas int

type Server struct {
	pb.UnimplementedChatServiceServer
	channel *amqp.Channel // Agregamos un campo para el canal de RabbitMQ
}

func Pedir_LLaves(cant_inicial int, cant_pedidas int)(int){

	if cant_inicial > 0 {
		if cant_pedidas == 0{
			num := int(cant_inicial/2)
			p := int(num/5)
			llaves_a_pedir := rand.Intn((num+p)-(num-p)+1) + (num - p)
			cant_llaves_pedidas=llaves_a_pedir
			return llaves_a_pedir
		}else{
			num := int(cant_inicial/2)
			p := int(num/5)
			llaves_a_pedir := rand.Intn((num+p)-(num-p)+1) + (num - p)
			cant_llaves_pedidas-=llaves_a_pedir
			return llaves_a_pedir
		}
	}else{
		return 0
	}
}


func (s *Server) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)

	// Enviamos un mensaje a RabbitMQ
	time.Sleep(5 * time.Second)
	inMessage:=string(in.Body)
	if inMessage == "LLaves Disponibles"{
		llaves_pedidas:=Pedir_LLaves(cant_registrados,cant_llaves_pedidas)
		err := s.channel.Publish(
			"",        // exchange
			"testing", // key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				//Body:        []byte(server_name+"-"+string(cant_llaves_pedidas)),
				Body:       []byte(server_name+"-"+strconv.Itoa(llaves_pedidas)),
				 // Enviamos el cuerpo del mensaje gRPC a RabbitMQ
			},
		)
		fmt.Println("Mande "+strconv.Itoa(llaves_pedidas)+" llaves")
		if err != nil {
			log.Printf("Error al publicar en RabbitMQ: %s", err)
		}
	}

	return &pb.Message{Body: "OK"}, nil
}

func(s *Server) SendKeys(ctx context.Context, in *pb.NumberRequest) (*pb.NumberResponse, error) {
	log.Printf("Receive Keys from client: %d", in.Number)
	cant_registrados-=int(in.Number)
	return &pb.NumberResponse{Response: "OK"}, nil

}

func main() {
	
	rand.Seed(time.Now().UnixNano())
	directorioActual, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        return
    }
    content, err := os.ReadFile(directorioActual+"/Regionales/Europa/parametros_de_inicio.txt")
	if err != nil {
		log.Fatal(err)
	}
	lineas := strings.Split(string(content), "\n")
	cant_registrados, err= strconv.Atoi(lineas[0])
	if err != nil {
		log.Fatal(err)
	}
	cant_llaves_pedidas=0
	
	
	server_name = "Europa"
	//addr_Rabbit := "localhost"
	addr_Rabbit := "dist106.inf.santiago.usm.cl"
	connection, err := amqp.Dial("amqp://guest:guest@" + addr_Rabbit + ":5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"testing",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	puerto := ":50054"
	lis, err := net.Listen("tcp", puerto)
	fmt.Printf("Escuchando %s\n", puerto)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{channel: channel} // Pasamos el canal de RabbitMQ al servidor gRPC
	pb.RegisterChatServiceServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)
}

package mq

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	connection    *amqp.Connection
	consumers     map[string]*amqp.Channel
	producers     map[string][]*amqp.Channel // Slice of channels for each queue
	queues        map[string]amqp.Queue
	producerChans map[string]chan string
	mutex         sync.Mutex
}

func NewRabbitMQService(connectionString string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{
		connection:    conn,
		consumers:     make(map[string]*amqp.Channel),
		producers:     make(map[string][]*amqp.Channel), // Changed to a slice
		queues:        make(map[string]amqp.Queue),
		producerChans: make(map[string]chan string),
	}, nil
}

func (r *RabbitMQService) DefineQueue(queueName string, numProducerThreads int, numConsumerThreads int, processFunc func(event string)) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.queues[queueName]; exists {
		log.Printf("Queue %s is already defined", queueName)
		return nil
	}

	// Declare queue
	channel, err := r.connection.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,  // Durable
		true,  // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	r.queues[queueName] = queue
	r.consumers[queueName] = channel
	r.producerChans[queueName] = make(chan string, 100) // Adjustable buffer size

	// Start producer threads
	fmt.Printf("Creating %d numProducerThreads\n", numProducerThreads)
	for i := 0; i < numProducerThreads; i++ {
		producerChannel, err := r.connection.Channel()
		if err != nil {
			log.Printf("Failed to create channel for producer %d: %v", i, err)
			continue
		}
		fmt.Printf("Created producer channel for %s\n", queueName)

		// Store each producer channel in a slice for the given queueName
		r.producers[queueName] = append(r.producers[queueName], producerChannel)

		// Start the producer goroutine for each channel
		go r.startProducer(queueName, r.producerChans[queueName], producerChannel)
	}

	// Start consumer threads
	fmt.Printf("Creating %d numConsumerThreads\n", numConsumerThreads)
	for i := 0; i < numConsumerThreads; i++ {
		go r.startConsumer(queueName, processFunc)
	}

	return nil
}

func (r *RabbitMQService) startProducer(queueName string, producerChan chan string, producerChannel *amqp.Channel) {
	for msg := range producerChan {
		err := producerChannel.Publish(
			"",        // Exchange
			queueName, // Routing key
			false,     // Mandatory
			false,     // Immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		}
	}
}

func (r *RabbitMQService) startConsumer(queueName string, processFunc func(event string)) {
	channel, err := r.connection.Channel()
	if err != nil {
		log.Printf("Failed to create channel for consumer: %v", err)
		return
	}
	fmt.Printf("Created consumer channel for %s\n", queueName)
	defer channel.Close()

	msgs, err := channel.Consume(
		queueName,
		"",    // Consumer tag
		false, // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	if err != nil {
		log.Printf("Failed to start consuming: %v", err)
		return
	}

	for msg := range msgs {
		processFunc(string(msg.Body))
		msg.Ack(false)
	}
}

func (r *RabbitMQService) SendMessage(queueName string, message string) error {
	if _, exists := r.producerChans[queueName]; !exists {
		return fmt.Errorf("Queue %s is not defined", queueName)
	}

	r.producerChans[queueName] <- message
	return nil
}

func (r *RabbitMQService) Close() {
	if r.connection != nil {
		r.connection.Close()
	}
	for _, channel := range r.consumers {
		channel.Close()
	}
	for _, channels := range r.producers {
		for _, channel := range channels {
			channel.Close()
		}
	}
}

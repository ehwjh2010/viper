package rabbitmq

import (
	"context"
	"github.com/ehwjh2010/viper/enums"
	"github.com/ehwjh2010/viper/helper/basic/collection"
	"github.com/ehwjh2010/viper/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type ProducerConf struct {
	Url      string
	Exchange Exchange
}

type RabbitProducer interface {
	Init() error
	SendMsg(ctx context.Context, key string, body []byte) error
	Close() error
}

type Producer struct {
	// 原生配置
	conf ProducerConf

	// 连接
	conn *amqp.Connection
	// 信道
	ch *amqp.Channel
	// 关闭通知信道
	closeNotifyChan chan *amqp.Error

	// 停止信道
	stopChan chan struct{}
	// 结束信道
	done chan struct{}
}

func NewProducer(conf ProducerConf) RabbitProducer {
	return &Producer{
		conf:     conf,
		stopChan: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// Watch 监听连接断开, 然后重连.
func (p *Producer) Watch() {
	oldConn := p.conn
	oldCh := p.ch

watchProducerLoop:
	for {
		select {
		case <-p.closeNotifyChan:
			if err := p.Setup(); err != nil {
				log.Errorf("rabbitmq producer reconnect failed, err: %s", err)
				time.Sleep(enums.FiveSecD)
			} else {
				oldCh.Close()
				oldConn.Close()
				oldConn, oldCh = p.conn, p.ch
				log.Infof("rabbitmq producer reconnect success")
			}
		case <-p.stopChan:
			break watchProducerLoop
		default:
			time.Sleep(enums.ThreeSecD)
		}
	}

	p.done <- struct{}{}
}

func (p *Producer) Setup() error {
	p.conf.Exchange.checkAndSet()

	conn, err := Connect(p.conf.Url)
	if err != nil {
		return err
	}

	// 获取信道
	ch, channelErr := conn.Channel()
	if channelErr != nil {
		return channelErr
	}

	// 声明交换机
	if err = ExchangeDeclare(ch, p.conf.Exchange); err != nil {
		return err
	}

	// 监听连接断开, 然后重连
	p.ch, p.conn = ch, conn
	p.closeNotifyChan = conn.NotifyClose(make(chan *amqp.Error))

	return nil
}

func (p *Producer) Init() error {
	err := p.Setup()
	if err != nil {
		return err
	}

	go p.Watch()
	return nil
}

// SendMsg 发送消息.
func (p *Producer) SendMsg(ctx context.Context, key string, body []byte) error {
	if collection.IsEmptyBytes(body) {
		return nil
	}

	err := p.ch.PublishWithContext(
		ctx,
		p.conf.Exchange.Name,
		key,
		false, // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false, // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain", // 消息内容的类型,
			Body:         body,
		})

	return err
}

// Close 关闭.
func (p *Producer) Close() error {
	p.stopChan <- struct{}{}
	<-p.done

	if err := p.ch.Close(); err != nil {
		log.Errorf("rabbitmq producer channel close failed, err: %s", err)
		return CancelChannelErr
	}

	if err := p.conn.Close(); err != nil {
		log.Errorf("rabbitmq producer connection close failed, err: %s", err)
		return CloseConnErr
	}

	log.Infof("rabbitmq producer close success")
	return nil
}

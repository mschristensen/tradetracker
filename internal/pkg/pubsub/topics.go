package pubsub

// Topic is a topic identifying a message stream.
type Topic string

// TradeTopic is the topic for trade messages.
var TradeTopic = Topic("trade")

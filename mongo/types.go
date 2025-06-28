package mongo

type PID struct {
	Key string `bson:"key"`
	Seq int    `bson:"seq"`
}

import asyncio
import signal
import sys
import time
import json
import pymongo
import redis
from nats.aio.client import Client as NATS
from stan.aio.client import Client as STAN

myclientmongo = pymongo.MongoClient("mongodb://35.237.232.19:27017/")
myclientredis = redis.Redis(
    host='35.237.232.19',
    port=6379, 
    password='',
    db=0
    )

async def run(loop):
    ##conexion con nats
    nc = NATS()
    await nc.connect(io_loop=loop, servers= ["nats://104.197.208.242:4222"])

    sc = STAN()
    await sc.connect("test-cluster","listener-3F45",nats=nc)

    async def cb(msg):
        print("Mensaje: (#{}): {}".format(msg.seq,msg.data))
        body = msg.data.decode("utf-8")
        data = json.loads(body)
        print(data["nombre"])
        print(data["departamento"])
        print("****")
        #inserto en mongo
        mydb = myclientmongo["proyecto"]
        mycol = mydb["casos"]
        x = mycol.insert_one(data)
        #inserto en redis
        myclientredis.set(str(x.inserted_id),body)

    subject = "lab"
    await sc.subscribe(subject,durable_name="3F45",cb = cb)

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    loop.run_forever()

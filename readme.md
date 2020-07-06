```
$ docker-compose up -d
```

MongoDBに集計対象、twitterのstream apiに対してtrackに指定するクエリデータを登録。optionsの項目に該当するtweetがor検索で返される
```
$ docker exec -it socialpool_mongo_1 /bin/bash
# mongo -u root -p example
> use ballots
> db.polls.insert({"title":"今の気分は?","options":["happy","sad","fail","win"]})
```

nsqのトピックのコンソール監視。publishから結構時間を要する。
```
$ docker exec -it socialpool_nsqd_1 sh
# nsq_tail --topic="votes" --lookupd-http-address=socialpool_nsqlookupd_1:4161
```
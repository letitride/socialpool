.envにtwitterにdeveloper登録したアプリの各tokenを記述しておく
```
SP_TWITTER_KEY=""
SP_TWITTER_SECRET=""
SP_TWITTER_ACCESSTOKEN=""
SP_TWITTER_ACCESSSECRET=""
```

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

curlを使ってapiをテストする
```
$ curl -XGET http://localhost:8080/polls/?Key=abc123
```
```
$ curl -XPOST http://localhost:8080/polls/?Key=abc123 -d '{"title": "調査のテスト", "options": ["one", "two", "three"]}'
```
```
$ curl -XDELETE http://localhost:8080/polls/5f071659f5bb22d8163cff21?Key=abc123
```
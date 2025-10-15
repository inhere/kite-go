
## simpleStore server

提供以下 API：

- `POST /filestore/{name}/add` 新增一行记录
- `GET /filestore/{name}/list` 获取所有记录

- `POST /jsonstore/{name}/add` 新增一条记录
- `GET /jsonstore/{name}/list` 获取所有记录

`POST /filestore/{name}/add` body:

```json
{
    "name": "xxx",
}
```

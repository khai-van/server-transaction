Command chạy

```bash
go get -u github.com/kataras/iris/v12@master

go mod tidy

go run main.go
```

## Cấu trúc của Project:

<u>/transaction</u>: sẽ chứa code thực hiện các hàm chức năng chính của repository (gồm tạo user, lấy user và chuyển tiền)

<u>/pkg</u>: sẽ chứa các package như database mongo, hàm dùng cho mutex lock dữ kiệu khi chuyển tiền 

<u>/api</u>: sẽ là folder chính để tạo server http

<u>/api/server.go</u>: tạo server http

<u>/api/configuration.go</u>: chứa config đọc từ file .yml về cấu hình http và database

<u>/api/router.go</u>: config đăng kí các router service

<u>/api/user</u>: chứa code chính về 1 service user router cho http server, folder này sẽ kết nối trực tiếp tới folder /transaction để sử dùng các hàm trong đây

## Các API

* \[GET]: /register?account=xZYXbQfBay 

response success:  

```json
{
  "msg": "Success",
  "payload": {
    "_id": "632d95c5b16df475ec49331a",
    "account": "xZYXbQfBay",
    "balance": 50
  }
}

```

* \[GET]: /transfer?to=xZYXbQfBay&from=nVrsWDeiTX&amount=16790 

response success:  

```json
{
  "msg": "Success",
  "payload": {
    "_id": "632d96be6100f624e2bb7c5b",
    "from": "nVrsWDeiTX",
    "to": "xZYXbQfBay",
    "amount": 16790,
    "created_at": 1663932094,
    "status": 1
  }
}
```

* \[GET]: /detail?account=xZYXbQfBay 

response success:  

```json
{
  "msg": "Success",
  "payload": {
    "transactions": [
      {
        "_id": "632d9600b16df475ec493435",
        "from": "xZYXbQfBay",
        "to": "nVrsWDeiTX",
        "amount": 50,
        "created_at": 1663931904,
        "status": 1
      },
      {
        "_id": "632d95fab16df475ec493422",
        "from": "nVrsWDeiTX",
        "to": "xZYXbQfBay",
        "amount": 100,
        "created_at": 1663931898,
        "status": 1
      }
    ],
    "user": {
      "_id": "632d95c5b16df475ec49331a",
      "account": "xZYXbQfBay",
      "balance": 50
    }
  }
}
```

## Các vấn đề đã giải quyết

<ul>
  <li>Partition Data: vì dữ liệu không nhiều các field hạn chế nên việc partition data là việc khó với project này. Tuy nhiên, em nghĩ việc partition data có thể làm được nếu ta xác định được field cần phải partition như, là các user trong 1 quốc gia,...  </li>
  <li>Data race: với 1 lượng lớn request transfer tới cùng 1 tài khoản, thì việc data race hoàn toàn có thể dẫn tới balance của user bị sai lệch. Và để tránh 1 request đợi quá lâu vì phải dùng cơ chế lock mặc định của golang(mỗi request tới lần lượt), thì trong project này em sử dụng 1 cơ chế lock dựa trên map (code ở /pkg/sync/waitmap.go). Với mỗi 1 transaction có from, to, gửi tới sẽ lock các transactions khác của 2 user này lại cho tới khi được unlock, giảm thiểu được các transaction không liên quan tới nhau bị lock, cải thiện được đáng kể performance nhưng vẫn giữ được tính chính sác</li>
  <li>Consistency: Vì hiện tại không có chia dữ liệu thành nhiều db nên có thể xem đây là strong consistency</li>

  <li>Nhược điểm:
    <ul>
        <li>Chưa có cơ chế hoàn lại balance nếu có sự cố khi update balance, chỉ có log và trace từ trạng thái của transaction</li>
        <li>Mã lỗi còn hạn chế</li>
        <li>Thiếu các cơ chế xác minh authentication</li>
    </ul>
  </li>
</ul>


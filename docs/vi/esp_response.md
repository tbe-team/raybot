# Định dạng phản hồi ESP

## 1. Định dạng tin nhắn

- Tất cả tin nhắn bắt đầu với ký tự `>`
- Tất cả tin nhắn kết thúc với `\r\n` (carriage return + line feed)
- Nội dung tin nhắn ở định dạng JSON

Ví dụ phản hồi:
```
sync: >{"type":0,"state_type":0;"data":{"speed":100}}\r\n
ack:  >{"type":1,"id":"abc123xyz";"status":0}\r\n
```

Cấu trúc JSON:
```json
{
  "type": <response_type>,
  "state_type": <state_type>, / "id":<id>,
  "data": <response_data>     / "status": <status>
}
```

### response_type

| Loại | Kiểu dữ liệu | Mô tả                     |
|------|--------------|---------------------------|
| 0    | uint8        | Đồng bộ trạng thái từ ESP |
| 1    | uint8        | ACK                       |

## 2. Phản hồi đồng bộ trạng thái (response_type = 0)

ESP tự động gửi phản hồi đồng bộ trạng thái đến ứng dụng.

Cấu trúc JSON:
```json
{
  "type": 0,
  "state_type": <state_type>,
  "data": {}
}
```

### state_type

| Loại | Kiểu dữ liệu | Mô tả                           |
|------|--------------|---------------------------------|
| a    | char         | Trạng thái cửa                  |
| b    | char         | Trạng thái động cơ đóng mở      |
| c    | char         | Trạng thái QR                   |
| d    | char         | Trạng thái cảm biến khoảng cách |

### data

Dữ liệu khác nhau cho mỗi loại trạng thái.

### 2.1. Trạng thái cửa

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| State | uint8 | state | Trạng thái cửa (0=Đóng/1=Mở) |


Ví dụ phản hồi:
```
>{"type":0,"state_type":"a","data":{"state":0}}\r\n
```

### 2.2. Trạng thái động cơ đóng mở

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| State | uint8 | state | Trạng thái của động cơ đóng mở (0=Đóng/1=Mở) |
| Speed | uint8 | speed | Tốc độ của động cơ di chuyển tính bằng % |
| IsRunning | uint8 | is_running | Động cơ có đang chạy hay không (0=false/1=true) |
| Enabled | uint8 | enabled | Cho phép động cơ di chuyển hoạt động hay không (0=false/1=true) |

Ví dụ phản hồi:
```
>{"type":0,"state_type":"b","data":{"state":0,"speed":50,"is_running":1,"enabled":1}}\r\n
```

### 2.3. Trạng thái QR Scanner

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| Code | string | code | Mã QR đọc được |

Ví dụ phản hồi:
```
>{"type":0,"state_type":"c","data":{"code":"abcxyz123"}}\r\n
```

### 2.4. Trạng thái cảm biến khoảng cách dưới

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| Front | uint8 | front | Khoảng cách đến vật phía trước tính bằng cm |
| Back | uint8 | back | Khoảng cách đến vật phía sau tính bằng cm |

Ví dụ phản hồi:
```
>{"type":0,"state_type":"d","data":{"front":100,"back":100}}\r\n
```

## 3. Phản hồi ACK

PIC gửi phản hồi ACK đến ứng dụng khi nhận được lệnh.

Cấu trúc JSON:
```json
{
  "type": 1,
  "id": <id>,
  "status": <status>
}
```

### id

- ID của lệnh
- Kiểu dữ liệu: string

### status

| Trường | Kiểu dữ liệu | Mô tả      |
|--------|--------------|------------|
| 0      | uint8        | Lỗi        |
| 1      | uint8        | Thành công |
# Định dạng phản hồi ESP

## 1. Định dạng tin nhắn

- Tất cả tin nhắn bắt đầu với ký tự `>`
- Tất cả tin nhắn kết thúc với `\r\n` (carriage return + line feed)
- Nội dung tin nhắn ở định dạng JSON

Ví dụ phản hồi:
```
>{"type":0,"state_type":0;"data":{"speed":100}}\r\n

```

Cấu trúc JSON:
```json
{
  "type": <response_type>,
  ...,
  "data": <response_data> 
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
| 0    | uint8        | Trạng thái cửa                  |
| 1    | uint8        | Trạng thái động cơ đóng mở      |
| 2    | uint8        | Trạng thái QR                   |
| 3    | uint8        | Trạng thái cảm biến khoảng cách |
| 4    | uint8        | Trạng thái pin                  |
| 5    | uint8        | Trạng thái sạc                  |
| 6    | uint8        | Trạng thái xả                   |

### data

Dữ liệu khác nhau cho mỗi loại trạng thái.

### 2.1. Trạng thái cửa

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả                               |
|--------|--------------|-----------|-------------------------------------|
| isOpen | bool         | is_open   | Trạng thái cửa (flase=Đóng/true=Mở) |


Ví dụ phản hồi:
```
>{"type":0,"state_type":0,"data":{"is_open":true}}\r\n
```

### 2.2. Trạng thái động cơ đóng mở

| Trường    | Kiểu dữ liệu | Khóa JSON  | Mô tả       |
|-----------|--------------|------------|-------------|
| State     | uint8        | state      | Trạng thái của động cơ đóng mở (0=Đóng/1=Mở) |
| Enabled   | uint8        | enabled    | Cho phép động cơ di chuyển hoạt động hay không (0=false/1=true) |
| Speed     | uint8        | speed      | Tốc độ của động cơ di chuyển tính bằng % |
| IsRunning | uint8        | is_running | Động cơ có đang chạy hay không (0=false/1=true) |


Ví dụ phản hồi:
```
>{"type":0,"state_type":1,"data":{"state":0,"enabled":1,"speed":50,"is_running":1}}\r\n
```

### 2.3. Trạng thái QR Scanner

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả          |
|--------|--------------|-----------|----------------|
| Code   | string       | code      | Mã QR đọc được |

Ví dụ phản hồi:
```
>{"type":0,"state_type":2,"data":{"code":"abcxyz123"}}\r\n
```

### 2.4. Trạng thái cảm biến khoảng cách dưới

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả       |
|--------|--------------|-----------|-------------|
| Under  | uint8        | under     | Khoảng cách đến vật phía dưới tính bằng cm |

Ví dụ phản hồi:
```
>{"type":0,"state_type":3,"data":{"under":100}}\r\n
```

### 2.5. Trạng thái pin

| Trường       | Kiểu dữ liệu | Khóa JSON     | Mô tả                              |
|--------------|--------------|---------------|------------------------------------|
| Current      | uint16       | current       | Dòng điện của pin tính bằng mA     |
| Temp         | uint8        | temp          | Nhiệt độ của pin tính bằng °C      |
| Voltage      | uint16       | voltage       | Điện áp của pin tính bằng mV       |
| CellVoltages | []uint16     | cell_voltages | Điện áp của từng cell tính bằng mV |
| Percent      | uint8        | percent       | Phần trăm pin                      |
| Fault        | uint8        | fault         | Trạng thái lỗi pin                 |
| Health       | uint8        | health        | Trạng thái sức khỏe pin            |

Ví dụ phản hồi:
```
>{"type":0,"state_type":0,"data":{"current":500,"temp":25,"voltage":12000,"cell_voltages":[4000,4000,4000,4000],"percent":80,"fault":0,"health":0}}\r\n
```

### 2.6. Trạng thái sạc

| Trường       | Kiểu dữ liệu | Khóa JSON     | Mô tả                               |
|--------------|--------------|---------------|-------------------------------------|
| CurrentLimit | uint16       | current_limit | Giới hạn dòng điện sạc tính bằng mA |
| Enabled      | uint8        | enabled       | Cho biết sạc có được bật hay không  |

Ví dụ phản hồi:
```
>{"type":0,"state_type":1,"data":{"current_limit":1000,"enabled":1}}\r\n
```

### 2.7. Trạng thái xả

| Trường       | Kiểu dữ liệu | Khóa JSON     | Mô tả                              |
|--------------|--------------|---------------|------------------------------------|
| CurrentLimit | uint16       | current_limit | Giới hạn dòng điện xả tính bằng mA |
| Enabled      | uint8        | enabled       | Cho biết xả có được bật hay không  |

Ví dụ phản hồi:
```
>{"type":0,"state_type":2,"data":{"current_limit":1000,"enabled":1}}\r\n
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
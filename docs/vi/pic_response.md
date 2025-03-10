# Định dạng phản hồi PIC

## 1. Định dạng tin nhắn

- Tất cả tin nhắn bắt đầu với ký tự `>`
- Tất cả tin nhắn kết thúc với `\r\n` (carriage return + line feed)
- Nội dung tin nhắn ở định dạng JSON

Ví dụ phản hồi:
```
>{"type":0,"data":{"current":500}}\r\n
```

Cấu trúc JSON:
```json
{
  "type": <response_type>,
  "data": <response_data>
}
```

### response_type

| Loại | Mô tả                     |
|------|---------------------------|
| 0    | Đồng bộ trạng thái từ PIC |
| 1    | ACK                       |

## 2. Phản hồi đồng bộ trạng thái (response_type = 0)

PIC tự động gửi phản hồi đồng bộ trạng thái đến ứng dụng.

Cấu trúc JSON:
```json
{
  "type": 0,
  "state_type": <state_type>,
  "data": {}
}
```

### state_type

| Loại | Mô tả                           |
|------|---------------------------------|
| 0    | Trạng thái pin                  |
| 1    | Trạng thái sạc                  |
| 2    | Trạng thái xả                   |
| 3    | Trạng thái cảm biến khoảng cách |
| 4    | Trạng thái động cơ nâng         |
| 5    | Trạng thái động cơ di chuyển    |

### data

Dữ liệu khác nhau cho mỗi loại trạng thái.

### 2.1. Trạng thái pin

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| Current | uint16 | current | Dòng điện của pin tính bằng mA |
| Temp | uint8 | temp | Nhiệt độ của pin tính bằng °C |
| Voltage | uint16 | voltage | Điện áp của pin tính bằng mV |
| CellVoltages | []uint16 | cell_voltages | Điện áp của từng cell tính bằng mV |
| Percent | uint8 | percent | Phần trăm pin |
| Fault | uint8 | fault | Trạng thái lỗi pin |
| Health | uint8 | health | Trạng thái sức khỏe pin |

Ví dụ phản hồi:
```
>{"type":0,"state_type":0,"data":{"current":500,"temp":25,"voltage":12000,"cell_voltages":[4000,4000,4000,4000],"percent":80,"fault":0,"health":0}}\r\n
```

### 2.2. Trạng thái sạc

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| CurrentLimit | uint16 | current_limit | Giới hạn dòng điện sạc tính bằng mA |
| Enabled | uint8 | enabled | Cho biết sạc có được bật hay không |

Ví dụ phản hồi:
```
>{"type":0,"state_type":1,"data":{"current_limit":1000,"enabled":1}}\r\n
```

### 2.3. Trạng thái xả

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| CurrentLimit | uint16 | current_limit | Giới hạn dòng điện xả tính bằng mA |
| Enabled | uint8 | enabled | Cho biết xả có được bật hay không |

Ví dụ phản hồi:
```
>{"type":0,"state_type":2,"data":{"current_limit":1000,"enabled":1}}\r\n
```

### 2.4. Trạng thái cảm biến khoảng cách

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| Front | uint16 | front | Khoảng cách đến vật phía trước tính bằng cm |
| Back | uint16 | back | Khoảng cách đến vật phía sau tính bằng cm |
| Down | uint16 | down | Khoảng cách đến vật bên dưới tính bằng cm |

Ví dụ phản hồi:
```
>{"type":0,"state_type":3,"data":{"front":100,"back":100,"down":100}}\r\n
```

### 2.5. Trạng thái động cơ nâng

| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| CurrentPosition | uint16 | current_position | Vị trí hiện tại của động cơ nâng (cm) |
| TargetPosition | uint16 | target_position | Vị trí mục tiêu của động cơ nâng (cm) |
| IsRunning | uint8 | is_running | Động cơ có đang chạy hay không (0=false/1=true) |
| Enabled | uint8 | enabled | Cho phép động cơ nâng hoạt động hay không (0=false/1=true) |

Ví dụ phản hồi:
```
>{"type":0,"state_type":4,"data":{"current_position":100,"target_position":200,"is_running":1,"enabled":1}}\r\n
```

### 2.6. Trạng thái động cơ di chuyển


| Trường | Kiểu dữ liệu | Khóa JSON | Mô tả |
|-------|-----------|----------|-------------|
| Direction | uint8 | direction | Hướng của động cơ di chuyển (0=Tiến/1=Lùi) |
| Speed | uint8 | speed | Tốc độ của động cơ di chuyển tính bằng % |
| IsRunning | uint8 | is_running | Động cơ có đang chạy hay không (0=false/1=true) |
| Enabled | uint8 | enabled | Cho phép động cơ di chuyển hoạt động hay không (0=false/1=true) |

Ví dụ phản hồi:
```
>{"type":0,"state_type":5,"data":{"direction":0,"speed":50,"is_running":1,"enabled":1}}\r\n
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

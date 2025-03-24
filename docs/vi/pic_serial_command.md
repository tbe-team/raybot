# Lệnh gửi đến PIC thông qua cổng serial

Tất cả lệnh được gửi đến PIC thông qua cổng serial đều có định dạng JSON.

Về phản hồi, PIC sẽ gửi phản hồi ACK đến ứng dụng sau khi nhận được lệnh và thực hiện lệnh thành công. Còn nếu lệnh thực hiện thất bại, PIC sẽ gửi lỗi đến ứng dụng.
Chi tiết: [PIC response ACK](pic_response.md#3-phản-hồi-ack)

## 1. Định dạng tin nhắn

- Tất cả tin nhắn bắt đầu với ký tự `>`
- Tất cả tin nhắn kết thúc với `\r\n` (carriage return + line feed)
- Nội dung tin nhắn ở định dạng JSON

Ví dụ lệnh:
```
>{"type":0,"data":{"current":500}}\r\n
```

Cấu trúc JSON:
```json
{
  "id": <id>,
  "type": <cmd_type>,
  "data": <cmd_data>
}
```

### id

- Kiểu dữ liệu: string
- Mô tả: ID của lệnh

### cmd_type

| Loại | Kiểu dữ liệu | Mô tả                      |
|------|--------------|----------------------------|
| 0    | uint8        | Cấu hình sạc pin           |
| 1    | uint8        | Cấu hình xả pin            |
| 2    | uint8        | Cấu hình động cơ nâng hạ   |
| 3    | uint8        | Cấu hình động cơ di chuyển |

### cmd_data

Tùy theo `cmd_type` mà `cmd_data` có nội dung khác nhau.

## 2. Cmd Data

### 2.1. Cấu hình sạc pin (cmd_type = 0)

| Tham số        | Kiểu dữ liệu | Mô tả                                |
|----------------|--------------|--------------------------------------|
| current_limit  | uint16       | Dòng điện tối đa                     |
| enable         | uint8        | Cho phép sạc pin (0: tắt, 1: bật)    |

Ví dụ:
```json
>{"id":"abc","type":0,"data":{"current_limit":123,"enable":1}}\r\n
```

### 2.2. Cấu hình xả pin (cmd_type = 1)

| Tham số        | Kiểu dữ liệu | Mô tả                                |
|----------------|--------------|--------------------------------------|
| current_limit  | uint16       | Dòng điện tối đa                     |
| enable         | uint8        | Cho phép xả pin (0: tắt, 1: bật)     |

Ví dụ:
```json
>{"id":"abc","type":1,"data":{"current_limit":123,"enable":1}}\r\n
```

### 2.3. Cấu hình động cơ nâng hạ (cmd_type = 2)

| Tham số         | Kiểu dữ liệu | Mô tả                                          |
|-----------------|--------------|------------------------------------------------|
| target_position | uint16       | Vị trí đích của động cơ nâng hạ                |
| enable          | uint8        | Cho phép động cơ nâng hạ chạy (0: tắt, 1: bật) |

Ví dụ:
```json
>{"id":"abc","type":2,"data":{"target_position":123,"enable":1}}\r\n
```

### 2.4. Cấu hình động cơ di chuyển (cmd_type = 3)

| Tham số         | Kiểu dữ liệu | Mô tả                                            |
|-----------------|--------------|--------------------------------------------------|
| direction       | uint8        | Hướng di chuyển (0: tiến, 1: lùi)                |
| speed           | uint8        | Tốc độ di chuyển (0-100)                         |
| enable          | uint8        | Cho phép động cơ di chuyển chạy (0: tắt, 1: bật) |

Ví dụ:
```json
>{"id":"abc","type":3,"data":{"direction":1,"speed":50,"enable":1}}\r\n
```

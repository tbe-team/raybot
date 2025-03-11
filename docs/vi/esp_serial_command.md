# Lệnh gửi đến ESP thông qua cổng serial

Tất cả lệnh được gửi đến ESP thông qua cổng serial đều có định dạng JSON.

Về phản hồi, ESP sẽ gửi phản hồi ACK đến ứng dụng sau khi nhận được lệnh và thực hiện lệnh thành công. Còn nếu lệnh thực hiện thất bại, PIC sẽ gửi lỗi đến ứng dụng.
Chi tiết: [ESP response ACK](esp_response.md#3-phản-hồi-ack)

## 1. Định dạng tin nhắn

- Tất cả tin nhắn bắt đầu với ký tự `>`
- Tất cả tin nhắn kết thúc với `\r\n` (carriage return + line feed)
- Nội dung tin nhắn ở định dạng JSON

Ví dụ lệnh:
```
>{"type":0,"data":{"speed":100}}\r\n
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

| Loại | Kiểu dữ liệu| Mô tả                      |
|------|-------------|----------------------------|
| 0    | uint8        | Cấu hình động cơ đóng mở   |

### cmd_data
Cấu hình động cơ đóng mở cửa
| Tham số         | Kiểu dữ liệu | Mô tả                                            |
|-----------------|--------------|--------------------------------------------------|
| state           | uint8        | Trạng thái cửa (0: đóng , 1: mở)                 |
| speed           | uint8        | Tốc độ di chuyển (0-100) %                        |
| enable          | uint8        | Cho phép động cơ đóng mở chạy (0: tắt, 1: bật)   |

Ví dụ:
```json
>{"id":"abc","type":0,"data":{"state":1,"speed":50,"enable":1}}\r\n
```
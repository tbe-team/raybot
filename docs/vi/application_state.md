# Trạng Thái Ứng Dụng

Chi tiết trạng thái robot mà ứng dụng quản lý.

## BatteryState (Trạng Thái Pin)
| Trường dữ liệu | Kiểu     | JSON Key      | Mô tả                                  |
|----------------|----------|---------------|----------------------------------------|
| Current        | uint16   | current       | Dòng điện của pin tính bằng mA         |
| Temp           | uint8    | temp          | Nhiệt độ của pin tính bằng °C          |
| Voltage        | uint16   | voltage       | Điện áp của pin tính bằng mV           |
| CellVoltages   | []uint16 | cell_voltages | Điện áp của từng cell tính bằng mV     |
| Percent        | uint8    | percent       | Phần trăm pin                          |
| Fault          | uint8    | fault         | Trạng thái lỗi pin                     |
| Health         | uint8    | health        | Trạng thái sức khỏe pin                |

## ChargeState (Trạng Thái Sạc)
| Trường dữ liệu | Kiểu    | JSON Key      | Mô tả                                    |
|----------------|---------|---------------|------------------------------------------|
| CurrentLimit   | uint16  | current_limit | Giới hạn dòng điện sạc tính bằng mA      |
| Enabled        | uint8   | enabled       | Cho phép sạc hay không (0=false/1=true)  |

## DischargeState (Trạng Thái Xả Pin)
| Trường dữ liệu | Kiểu    | JSON Key      | Mô tả                                      |
|----------------|---------|---------------|--------------------------------------------|
| CurrentLimit   | uint16  | current_limit | Giới hạn dòng điện xả tính bằng mA         |
| Enabled        | uint8   | enabled       | Cho phép xả pin hay không (0=false/1=true) |

## DistanceSensorState (Trạng Thái Cảm Biến Khoảng Cách)
| Trường dữ liệu | Kiểu    | JSON Key       | Mô tả                                               |
|----------------|---------|----------------|-----------------------------------------------------|
| FrontDistance  | uint16  | front_distance | Khoảng cách đến phía trước robot (cm)               |
| BackDistance   | uint16  | back_distance  | Khoảng cách đến phía sau robot (cm)                 |
| DownDistance   | uint16  | down_distance  | Khoảng cách từ đáy robot đến nắp hộp chứa hàng (cm) |

## LiftMotorState (Trạng Thái Động Cơ Nâng)
| Trường dữ liệu   | Kiểu    | JSON Key         | Mô tả                                                      |
|------------------|---------|------------------|------------------------------------------------------------|
| CurrentPosition  | uint16  | current_position | Vị trí hiện tại của động cơ nâng (cm)                      |
| TargetPosition   | uint16  | target_position  | Vị trí mục tiêu của động cơ nâng (cm)                      |
| IsRunning        | uint8   | is_running       | Động cơ có đang chạy hay không (0=false/1=true)            |
| Enabled          | uint8   | enabled          | Cho phép động cơ nâng hoạt động hay không (0=false/1=true) |

## DriveMotorState (Trạng Thái Động Cơ Di Chuyển)
| Trường dữ liệu | Kiểu                | JSON Key      | Mô tả                                                           |
|----------------|---------------------|---------------|-----------------------------------------------------------------|
| Direction      | uint8               | direction     | Hướng của động cơ di chuyển (0=Tiến/1=Lùi)                      |
| Speed          | uint8               | speed         | Tốc độ của động cơ di chuyển (%)                                |
| IsRunning      | uint8               | is_running    | Động cơ có đang chạy hay không (0=false/1=true)                 |
| Enabled        | uint8               | enabled       | Cho phép động cơ di chuyển hoạt động hay không (0=false/1=true) |

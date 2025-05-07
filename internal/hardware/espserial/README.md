### ✅ `Client.Read()` Test Cases

| #   | Input (mock serial data)           | Expected Output                           |
|-----|------------------------------------|-------------------------------------------|
| 1   | `>data\r\n`                        | `data`                    
| 2   | `>>\rdata\r\n`                     | `>\rdata`                    
| 3   | `>data\x00\r\n`                    | `data\x00`               
| 4   | `>>data\x00\r\n`                   | `>data\x00`          
| 5   | `>>data\x00\r\n\r\n`               | `>data\x00`      
| 6   | `random>message here\r\ngarbage`   | `message here`             
| 7   | `>null\r\n\x00\x00`                | `null`                    
| 7   | `>\r\n`                            | ``          
| 8   | *(no port connected)*              | `ErrESPSerialNotConnected` 
| 9   | `>data without end`                | EOF               
| 10  | `data only`                        | EOF                
| 11  | *(context cancelled)*              | Context cancelled      

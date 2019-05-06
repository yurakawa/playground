- build

`make build`

- exec

```
### 暗号鍵作成

$ ./build/keygen
9YsltIdBuMHJgcO1O+e47mBaHWId1/Sv6bZpLJ3BsxY=

$ export ANGO_KEY=9YsltIdBuMHJgcO1O+e47mBaHWId1/Sv6bZpLJ3BsxY=

### 暗号化

cat secret.txt | ./build/ango > encrypted_secret.txt

### 復号

cat encrypted_secret.txt | ./build/fukugo
```
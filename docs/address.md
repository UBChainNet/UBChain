# 地址生成 文档
## 目录

### 工具包
```
github.com/jhdriver/UBChain/param
github.com/jhdriver/UBChain/ut
github.com/jhdriver/UBChain/common/hasharry
```


###  生成BIP39助记词

```
e, _ := ut.Entropy()
m, _ := ut.Mnemonic(e)
```

### 生成secp256k1私钥

```
key, _ := ut.MnemonicToEc(m)
```

###  生成地址

```
addr, _ := ut.GenerateUBCAddress(param.TestNet, key.PubKey())
```

### 校验地址
```
ut.CheckUBCAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy")
```

### 生成token地址

```
ut.GenerateTokenAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy", "HFC")
```

### 校验token地址

```
ut.CheckTokenAddress(param.TestNet, "UbQQhJ4zmp4wLQ4Li6tm7zigopaeGrWxSvy")
```
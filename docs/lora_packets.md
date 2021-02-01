
# Join Request

Message Type = Join Request
AppEUI = 70B3D57ED00000DC
DevEUI = 00AFEE7CF5ED6F1E
DevNonce = CC85
MIC = 587FE913

=>  00DC0000D07ED5B3701E6FEDF57CEEAF0085CC587FE913


# Join Accept decoding



example:
Payload : 204DD85AE608B87FC4889970B7D2042C9E72959B0057AED6094B16003DF12DE145
Application Key: B6B53F4A168A7A88BDF7EA135CE9CFCA

Message Type = Join Accept
AppNonce = E5063A
NetId = 000013
DevAddr = 26012E43
DLSettings = 03
RXDelay = 01
CFList = 184F84E85684B85E84886684586E8400
        (8671000, 8673000, 8675000, 8677000, 8679000)
MIC = 55121DE0
NB: The network server uses an AES decrypt operation in ECB mode to encrypt the join-accept
message so that the end-device can use an AES encrypt operation to decrypt the message.
This way an end-device only has to implement AES encrypt but not AES decrypt.

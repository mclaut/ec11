<h1>EC11</h1>
Golang package for EC-11 encoder, based on periph/io

<i>New(DTpin1,CLKpin2,SWpin3 int) (EncoderT, error)</i> - create Encoder object<br>
<i>Start() chan int8</i> - start & return chan int8<br>
Package send to chan int8:<br>
-1 - 1 rotate counter clockwise<br>
1 - 1 rotate clockwise<br>
0 - pressed button<br>

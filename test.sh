 #!/bin/bash 

for i in `seq 1 100`;
do
 echo "127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 1234" >> test.log
done

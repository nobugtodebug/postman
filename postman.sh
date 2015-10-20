 echo "email,Reason" > ad4tangseng/send_buf.csv
 head -n 7000 ad4tangseng/maillist.csv | tail -n 500 >> ad4tangseng/send_buf.csv
 ./postman -server smtp.exmail.qq.com  \
     -port 25 \
     -sender ts99@tangseng99.com \
     -user ts99@tangseng99.com \
     -password tangseng99 \
     -subject "轻松搞定不能访问国外网站的问题" \
     -text ad4tangseng/ts99.txt \
     -html ad4tangseng/ts99.html \
     -csv ad4tangseng/send_buf.csv  \
     -rand ad4tangseng/randDesc.json \
     -fmin 5 \
     -freq 2

# echo "email,Reason" > ad4tangseng/send_buf.csv
# head -n 4000 ad4tangseng/maillist.csv | tail -n 500 >> ad4tangseng/send_buf.csv
# ./postman -server smtp.exmail.qq.com  \
#     -port 25 \
#     -sender noreply@tangseng99.com \
#     -user noreply@tangseng99.com \
#     -password tangseng99 \
#     -subject "轻松搞定不能访问国外网站的问题" \
#     -text ad4tangseng/ts99.txt \
#     -html ad4tangseng/ts99.html \
#     -csv ad4tangseng/send_buf.csv  \
#     -rand ad4tangseng/randDesc.json \
#     -fmin 5 \
#     -freq 2

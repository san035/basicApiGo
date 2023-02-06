$pathScript = $MyInvocation.MyCommand.Path | split-path -parent
cd $pathScript\cmd\
write-host `n"Path script: "$pathScript

write-host `n"Project Compilation"
gox -osarch="linux/amd64" -output="../../app_basic" -tags go_tarantool_ssl_disable
cd ..

write-host `n"Stopping the application on the server"
ssh askorohodov@bpm.dev.itkn.ru 'tmux kill-session -t profile'
ssh askorohodov@bpm.dev.itkn.ru '/opt/BackEnd/api-service/app_basic down'

write-host `n"Copy to server"
scp C:/go/api-service/app_basic askorohodov@bpm.dev.itkn.ru:/opt/BackEnd/api-service/app_basic

write-host `n"Running on the server"
ssh -t askorohodov@bpm.dev.itkn.ru 'tmux new -s profile /opt/BackEnd/api-service/app_basic'

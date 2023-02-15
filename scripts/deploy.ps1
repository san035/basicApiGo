$curFolder=$pwd
echo `n"Currenr folder : "$curFolder

$AppName = "app_basic"
$FolderAppName = "api-service/"

$pathScript = $MyInvocation.MyCommand.Path | split-path -parent
write-host `n"Path script: "$pathScript

if ( 1 )
{
write-host `n"1. Project Compilation"
cd $pathScript\..\cmd\
gox -osarch="linux/amd64" -output="../"$AppName -tags go_tarantool_ssl_disable
cd ..
}

write-host `n"2. Stopping the application on the server"
ssh askorohodov@bpm.dev.itkn.ru 'tmux kill-session -t '$AppName
$cmdDown="/opt/BackEnd/"+$FolderAppName+$AppName+" down"
ssh askorohodov@bpm.dev.itkn.ru $cmdDown

write-host `n"3. Copy to server"
$dist="askorohodov@bpm.dev.itkn.ru:/opt/BackEnd/"+$FolderAppName+$AppName
scp $AppName $dist

#write-host `n"4. set right"
#$distRight="sudo chmod ugo+rwx /opt/BackEnd/"+$FolderAppName+$AppName
#ssh -tt askorohodov@bpm.dev.itkn.ru $distRight

write-host `n"4. Running on the server"
$cmdTmux="tmux new -s "+ $AppName + " /opt/BackEnd/"+$FolderAppName+$AppName
ssh -tt askorohodov@bpm.dev.itkn.ru $cmdTmux

cd $curFolder
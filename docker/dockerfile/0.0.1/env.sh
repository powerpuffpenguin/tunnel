DockerVarName="test-tunnel-0.0.1"
DockerVarShell="sh"
function before_build
{
    rm root -rf
    mkdir root/opt/tunnel -p

    cp "$ProjectDir/../bin/tunnel" root/opt/tunnel
    cp "$ProjectDir/../bin/etc" root/opt/tunnel/ -r
}

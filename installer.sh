#!/usr/bin/env bash 

check_dependencies() {
    if [[ $(which go) == "" ]]; then 
        echo -e "[+] Go not found\n[+] Installing Go"
        sudo apt-get install golang -y 
        if [[ $(echo $?) != 0 ]]; then 
            echo "[!] Error trying install Go"
            exit 1;
        fi
    else
        creating_directories
    fi
}


creating_directories() {
    paths=( "/etc/systemd/system/" "/etc/systemActivator/" )
    for path in "${paths[@]}"; do 
        if [[ ! -d $path ]]; then 
            echo "[+] Creating directory $path"
            sudo mkdir -p $path 
        fi 
    done
    compile_programm
}

compile_programm() {
    go build -o systemActivator 
    if [[ $(echo $?) == 0 ]]; then 
        echo "[+] Moving binary"
        sudo mv systemActivator /usr/local/bin/
        create_systemd_file
    else 
        echo "[!] Error compiling programm"

    fi
}


create_systemd_file() {
    echo -e "[+] Creating .service file"
    cat << EOF | sudo tee /etc/systemd/system/systemActivator.service > /dev/null
[Unit]
Description=System Activator bot 
After=network-online.target 
Wants=network-online.target 

[Service]
EnvironmentFile=/etc/systemActivator/.env 
ExecStart=/usr/local/bin/systemActivator
User=root
Group=root 

[Install]
WantedBy=multi-user.target 
EOF
    if [[ $(echo $?) != 0 ]]; then 
        echo -e "[!] Error creating .service file" 
        exit 1;
    fi
    create_env_file
}

create_env_file() {
    echo "[+] Creating .env files"
    read -sp "[+] Please paste your Telegram bot Token: " telegramToken
    echo ""
    read -sp "[+] Please past your Telegram Chat ID: " telegramChatId
    echo ""
    printf "TOKEN_ID=%s\nCHAT_ID=%s\n" "$telegramToken" "$telegramChatId" | sudo tee /etc/systemActivator/.env > /dev/null
    if [[ $(echo $?) != 0 ]]; then 
        echo "[!] Error creating .env"
        exit 1;
    fi
    restart_services
}

restart_services() {
    sudo systemctl daemon-reload 
    sudo systemctl start systemActivator
    sudo systemctl enable systemActivator
}

main() {
    check_dependencies
}

main


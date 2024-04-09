provider "aws" {
  region     = "us-east-1"
  access_key = ""
  secret_key = ""
}

resource "aws_security_group" "ec2_sg" {
  name = "aws_security_group_ec2_sg"
  description = "aws security group ec2 sg"

  ingress {
    description = "search-client access"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "ssh access"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "search-api access"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 4080
    to_port     = 4080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  user_data = <<-EOF
    sudo apt update 
    sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update 
    apt-cache policy docker-ce
    sudo apt install -y docker-ce

    wget https://golang.org/dl/go1.21.4.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin

    git clone https://github.com/heisjuanda/Go-and-Vue-journey.git /home/ec2-user/truora
    sudo curl -L https://github.com/zinclabs/zinc/releases/download/v0.3.6/zinc_0.3.6_Linux_x86_64.tar.gz -o /home/ec2-user/truora/Back/zinc.tar.gz
    sudo tar -xf /home/ec2-user/truora/Back/zinc.tar.gz -C /home/ec2-user/truora/Back/
    cd /home/ec2-user/truora/Back
    sudo mkdir data
    cd /home/ec2-user
    sudo wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz
    tar -xvzf enron_mail_20110402.tgz
  EOF

  tags = {
    Name = "web-server"
  }
}

output "public_ip" {
  value = aws_instance.web.public_ip
}

resource "aws_lb" "this" {
  name               = "${var.project}-nlb"
  load_balancer_type = "network"
  internal           = false
  subnets            = [var.public_subnet_id]
  security_groups    = [aws_security_group.nlb.id]
  tags = {
    Name = "${var.project}-nlb"
  }
}

resource "aws_lb_target_group" "traefik" {
  name        = "${var.project}-traefik"
  port        = 80
  protocol    = "TCP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    protocol = "TCP"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.traefik.arn
  }
}

resource "aws_security_group" "nlb" {
  name        = "${var.project}-nlb-sg"
  description = "Security group for ingress controller"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Allow HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-ingress-controller-sg"
  }
}

#!/bin/bash
# ec2.sh — Manually start or stop the EC2 instance
# Usage:
#   ./scripts/ec2.sh stop [environment]
#   ./scripts/ec2.sh start [environment]
#   ./scripts/ec2.sh status [environment]

set -e

ACTION=${1:-status}
ENV=${2:-shared}
REGION="ap-southeast-1"
PROJECT="openclaw"
INSTANCE_NAME="${PROJECT}-${ENV}-${PROJECT}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

get_instance_id() {
  aws ec2 describe-instances \
    --filters "Name=tag:Name,Values=${INSTANCE_NAME}" "Name=instance-state-name,Values=pending,running,shutting-down,stopping,stopped" \
    --region "$REGION" \
    --query "Reservations[0].Instances[0].InstanceId" \
    --output text 2>/dev/null | grep -v "None" || echo ""
}

get_status() {
  local id=$1
  if [ -z "$id" ]; then
    echo "not-found"
    return
  fi
  aws ec2 describe-instances \
    --instance-ids "$id" \
    --region "$REGION" \
    --query "Reservations[0].Instances[0].State.Name" \
    --output text 2>/dev/null || echo "not-found"
}

print_status() {
  local id=$1
  STATUS=$(get_status "$id")
  case "$STATUS" in
    running)    echo -e "  Status: ${GREEN}● running${NC}" ;;
    stopped)    echo -e "  Status: ${RED}● stopped${NC}" ;;
    stopping)   echo -e "  Status: ${YELLOW}⏳ stopping...${NC}" ;;
    pending)    echo -e "  Status: ${YELLOW}⏳ pending (starting)...${NC}" ;;
    not-found)  echo -e "  Status: ${RED}✗ not found${NC}" ;;
    *)          echo -e "  Status: ${BLUE}${STATUS}${NC}" ;;
  esac
}

echo ""
echo "🖥️  EC2: ${INSTANCE_NAME} (${REGION})"
echo "──────────────────────────────────────────"

INSTANCE_ID=$(get_instance_id)

if [ -z "$INSTANCE_ID" ] && [ "$ACTION" != "help" ]; then
  echo -e "  Status: ${RED}✗ not found${NC}"
  echo ""
  exit 1
fi

case "$ACTION" in
  stop)
    STATUS=$(get_status "$INSTANCE_ID")
    if [ "$STATUS" = "stopped" ]; then
      echo -e "${YELLOW}Already stopped.${NC}"
      print_status "$INSTANCE_ID"
    elif [ "$STATUS" = "running" ]; then
      echo "Stopping EC2 instance ($INSTANCE_ID)..."
      aws ec2 stop-instances \
        --instance-ids "$INSTANCE_ID" \
        --region "$REGION" > /dev/null
      echo -e "${GREEN}✓ Stop initiated.${NC} (takes ~1-2 minutes)"
      print_status "$INSTANCE_ID"
    else
      echo -e "${RED}Cannot stop — current status: ${STATUS}${NC}"
      exit 1
    fi
    ;;

  start)
    STATUS=$(get_status "$INSTANCE_ID")
    if [ "$STATUS" = "running" ]; then
      echo -e "${YELLOW}Already running.${NC}"
      print_status "$INSTANCE_ID"
    elif [ "$STATUS" = "stopped" ]; then
      echo "Starting EC2 instance ($INSTANCE_ID)..."
      aws ec2 start-instances \
        --instance-ids "$INSTANCE_ID" \
        --region "$REGION" > /dev/null
      echo -e "${GREEN}✓ Start initiated.${NC} (takes ~1-2 minutes to become available)"
      print_status "$INSTANCE_ID"
    else
      echo -e "${RED}Cannot start — current status: ${STATUS}${NC}"
      exit 1
    fi
    ;;

  status)
    echo -e "  Instance ID: $INSTANCE_ID"
    print_status "$INSTANCE_ID"
    ;;

  *)
    echo -e "${RED}Unknown action: ${ACTION}${NC}"
    echo ""
    echo "Usage: $0 [stop|start|status] [environment]"
    echo "  environment defaults to 'shared'"
    exit 1
    ;;
esac

echo ""

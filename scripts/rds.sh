#!/bin/bash
# rds.sh — Manually start or stop the RDS instance
# Usage:
#   ./scripts/rds.sh stop [environment]
#   ./scripts/rds.sh start [environment]
#   ./scripts/rds.sh status [environment]

set -e

ACTION=${1:-status}
ENV=${2:-dev}
REGION="ap-southeast-1"
DB_IDENTIFIER="dora-metrics-${ENV}-db"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

get_status() {
  aws rds describe-db-instances \
    --db-instance-identifier "$DB_IDENTIFIER" \
    --region "$REGION" \
    --query 'DBInstances[0].DBInstanceStatus' \
    --output text 2>/dev/null || echo "not-found"
}

print_status() {
  STATUS=$(get_status)
  case "$STATUS" in
    available)  echo -e "  Status: ${GREEN}● available${NC}" ;;
    stopped)    echo -e "  Status: ${RED}● stopped${NC}" ;;
    stopping)   echo -e "  Status: ${YELLOW}⏳ stopping...${NC}" ;;
    starting)   echo -e "  Status: ${YELLOW}⏳ starting...${NC}" ;;
    not-found)  echo -e "  Status: ${RED}✗ not found${NC}" ;;
    *)          echo -e "  Status: ${BLUE}${STATUS}${NC}" ;;
  esac
}

echo ""
echo "🗄️  RDS: ${DB_IDENTIFIER} (${REGION})"
echo "──────────────────────────────────────────"

case "$ACTION" in
  stop)
    STATUS=$(get_status)
    if [ "$STATUS" = "stopped" ]; then
      echo -e "${YELLOW}Already stopped.${NC}"
      print_status
    elif [ "$STATUS" = "available" ]; then
      echo "Stopping RDS instance..."
      aws rds stop-db-instance \
        --db-instance-identifier "$DB_IDENTIFIER" \
        --region "$REGION" > /dev/null
      echo -e "${GREEN}✓ Stop initiated.${NC} (takes ~1-2 minutes)"
      print_status
    else
      echo -e "${RED}Cannot stop — current status: ${STATUS}${NC}"
      exit 1
    fi
    ;;

  start)
    STATUS=$(get_status)
    if [ "$STATUS" = "available" ]; then
      echo -e "${YELLOW}Already running.${NC}"
      print_status
    elif [ "$STATUS" = "stopped" ]; then
      echo "Starting RDS instance..."
      aws rds start-db-instance \
        --db-instance-identifier "$DB_IDENTIFIER" \
        --region "$REGION" > /dev/null
      echo -e "${GREEN}✓ Start initiated.${NC} (takes ~3-5 minutes to become available)"
      print_status
    else
      echo -e "${RED}Cannot start — current status: ${STATUS}${NC}"
      exit 1
    fi
    ;;

  status)
    print_status
    ;;

  *)
    echo -e "${RED}Unknown action: ${ACTION}${NC}"
    echo ""
    echo "Usage: $0 [stop|start|status] [environment]"
    echo "  environment defaults to 'dev'"
    exit 1
    ;;
esac

echo ""

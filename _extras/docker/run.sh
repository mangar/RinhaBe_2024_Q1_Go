

docker run -d --name api00 -p 8000:8000 \
    -e FLAG_DEBUG="True" \
    -e SERVER_NAME="SERVER_API00" \
    -e LOG_OUTPUT_DIR="./" \
    -e DB_CONNECTION="postgres://postgres:password@192.168.1.110:5432/postgres?sslmode=disable" \
 mangar/rinhabe_2024_q1_go:0.0.1


# docker run -d --name api00 -p 8000:8000 \
#     -e FLAG_DEBUG="True" \
#     -e SERVER_NAME="SERVER_API00" \
#     -e LOG_OUTPUT_DIR="/" \
#     -e DB_CONNECTION="postgres://postgres:password@192.168.1.110:5432/postgres" \
#     -e DB_POOL_SIZE="30" \
#     -e DB_MAX_OVERFLOW="40" \
#     -e DB_POOL_TIMEOUT="10" \
#     -e DB_POOL_RECYCLE="600" \
#  mangar/rinhabe_2024_q1_go:0.0.1

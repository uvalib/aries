# run application

# run from here, since application expects web template in web/ and writes pdfs to tmp/
cd bin; ./aries -redis_host "$ARIES_REDIS_HOST" -redis_port "$ARIES_REDIS_PORT" -redis_db "$ARIES_REDIS_DB" -redis_pass "$ARIES_REDIS_PASS" -redis_prefix "$ARIES_REDIS_PREFIX"

#
# end of file
#

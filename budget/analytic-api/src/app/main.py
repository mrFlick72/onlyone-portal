import logging
import logging.handlers
import uvicorn
import os

log_file = os.getenv("LOG_FILE_LOCATION", "logs/app.log")
log_level_str = os.getenv("LOG_LEVEL", "INFO").upper()

level = getattr(logging, log_level_str, logging.INFO)

root = logging.getLogger()
root.setLevel(level)
root.handlers.clear()

fmt = logging.Formatter("%(asctime)s %(name)s %(levelname)s: %(message)s")

ch = logging.StreamHandler()
ch.setLevel(level)
ch.setFormatter(fmt)
root.addHandler(ch)

os.makedirs(os.path.dirname(log_file) or ".", exist_ok=True)
fh = logging.handlers.RotatingFileHandler(log_file, maxBytes=10_000_000, backupCount=5)
fh.setLevel(level)
fh.setFormatter(fmt)
root.addHandler(fh)

def main() -> None:
    root.info("Starting server with log level: %s", log_level_str)
    uvicorn.run("app.server:app", host="0.0.0.0", port=3035, reload=False)


if __name__ == "__main__":
    main()

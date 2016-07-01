import os, sys
import globalConfig
# def script_path():
#     path = os.path.realpath(sys.path[0])
#     if os.path.isfile(path):
#         path = os.path.dirname(path)
#     return os.path.abspath(path)

# g_current_dir = script_path()

def main():
	print globalConfig.g_current_dir
	tpcmd = "redis-server " + os.path.join(globalConfig.g_current_dir, "redis.windows.conf")
	os.system(tpcmd)


if __name__ == "__main__":
	main()

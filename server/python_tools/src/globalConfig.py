#coding=utf-8
#!/usr/bin/python

"""
config 
"""
import os,sys



#common config
def script_path():
    path = os.path.realpath(sys.path[0])
    if os.path.isfile(path):
        path = os.path.dirname(path)
    return os.path.abspath(path)



g_current_dir = script_path()
g_project_root_dir = os.path.join(g_current_dir, "..", "..")

g_service_dir = os.path.join(g_project_root_dir, "service")

# #g_python_config_dir = os.path.join(g_project_root_dir, "pythonConfig") 
# #sys.path.append(g_python_config_dir)
# #import branch_config
# #branch_config.g_branch_table_res_dir

# g_cache_dir = os.path.join(g_project_root_dir, "..", "cacheClient")
# g_cache_dir2 = os.path.join(g_project_root_dir, "..", "cacheClient2")
# g_share_dir = os.path.join(g_project_root_dir, "branch_config")
# g_server_test_dir = os.path.join(g_project_root_dir, "..", "hall")
# g_server_test_dir2 = os.path.join(g_project_root_dir, "..", "svnData2")

# g_ui_dir = os.path.join(g_project_root_dir, "resources", "ChessCard", "CocosStudio")
# g_model_res_dir = os.path.join(g_project_root_dir, "resources", "model_res")
# g_image_dir = os.path.join(g_project_root_dir, "resources","image_res")
# g_local_dir = os.path.join(g_project_root_dir, "resources","local")
# g_sound_dir = os.path.join(g_project_root_dir, "resources","sound")
# #g_game_tables_dir = os.path.join(g_project_root_dir, "resources","table_res")

# g_games_dir = os.path.join(g_project_root_dir, "hall","xm","games")
# g_games_res_dir = os.path.join(g_project_root_dir, "hall","res")
# g_games_common_dir = os.path.join(g_games_dir, "common")
# g_hall_table_dir = os.path.join(g_games_dir, "hall","config")
# #g_android_framework_dir = os.path.join(g_project_root_dir, "hall","frameworks","runtime-src","proj.android_apus")
# g_android_framework_dir = os.path.join(g_project_root_dir, "hall","frameworks","runtime-src","proj.android_no_anysdk")

# g_common_framework_dir = os.path.join(g_project_root_dir, "hall","xm","fw")

# g_game_tables_dir = os.path.join(g_share_dir, "table_res")
# g_game_tables_common = os.path.join(g_share_dir, "table_common")
# g_sensitive_word_table = os.path.join(g_share_dir, "sensitive_word", "sensitive_word.xlsx") 

# g_artist_dir = os.path.join(g_project_root_dir, "..","DT","美术","印度大厅")

# g_unique_res_dir = os.path.join(g_project_root_dir, "resources_unique")
# g_unique_asset_dir = os.path.join(g_unique_res_dir, "ChessCard", "CocosStudio", "assets")
# g_unique_img_dir = os.path.join(g_unique_res_dir, "image_res")
# g_unique_model_res_dir = os.path.join(g_unique_res_dir, "model_res")

# g_pack_image_temp_dir = os.path.join(g_cache_dir, "image_res_temp")
# g_compress_tool_dir = os.path.join(g_current_dir, "compressTool")

# g_py_config_dir = os.path.join(g_project_root_dir, "hall", "pyconfig")

# sys.path.append(g_py_config_dir)
# sys.path.append(os.path.join(g_py_config_dir, "table_config"))

SERVICES_LIST = ["cfgcenter",
"oplogs",
"toplist",
"recharge",
"coredata_data",
"coredata_logic",
"statuskeeper",
"auth",
"hall",
"version",
"gateway",
"gamelist",
"matchsign",
"gamegw",
"social",
"gamemaster",
"minigames"]

if __name__ == "__main__":
    print g_current_dir

    
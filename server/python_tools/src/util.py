#!/usr/bin/python
#coding=utf-8

"""
工具类
"""	
import shutil
import os,sys
#from PIL import Image
#gogogogo!
import filecmp
import globalConfig

def check_file(file):
	filename = os.path.split(file)[1]
	if filename.endswith(".ExportJson") or filename.endswith(".Json"):
		return 
	if file.endswith(".DS_Store"):
		return 
	if not filename.islower():
		print "\n***********\nFuck!! file %s, has upper letter, change it quickly!!!\n************\n"%file
		sys.exit(1)

def is_same_file(src_file, des_file):
	return filecmp.cmp(src_file, des_file, 0)


def check_lower_dir(arg,dirname,names):
	for name in names:
		file = os.path.join(dirname, name)
		if not os.path.isfile(file):
			continue
		check_file(file)		


def copy_file_to_dir(file, to_path, is_check=True):
	if file.endswith(".DS_Store"):
		return 	
	if is_check:
		check_file(file)
	try:
		os.makedirs(to_path)
	except:
		pass
	shutil.copy2(file, to_path)
	
	
def copy_dir_to_dir(from_path, to_path, is_check=True):
	if  os.path.isdir(to_path):
		delete_dir(to_path)
	if is_check:
		os.path.walk(from_path,check_lower_dir,())
	try:
		os.makedirs(to_path)
	except:
		pass
	shutil.copytree(from_path, to_path)
	

def delete_dir(dir):
	print "delete dir: ",dir
	try:
		shutil.rmtree(dir)	
	except os.error, err:
		print err

def delete_file(file):
	try:
		os.remove(file)
		print "delete file: ", file
	except:
		print file, "is not exist"
	
def	recreate_dir(dir):
	if os.path.isdir(dir):
		delete_dir(dir)
	print "recreate_dir: ",dir
	try:
		os.makedirs(dir)
	except os.error, err:
		print err

		
def is_endswith(s, items):
	for item in items:
		if s.endswith(item):
			return True
	return False

def is_startswith(s, items):
	for item in items:
		if s.startswith(item):
			return True
	return False
		
def clear_despath(des_path, endWithItems, exeludeItems):
	for eachFile in os.listdir(des_path):
		currentFilePath = os.path.join(des_path, eachFile)
		if endWithItems != None:
			if not is_endswith(eachFile, endWithItems):
				continue
		if exeludeItems != None:
			if is_endswith(eachFile, exeludeItems):
				continue
		if os.path.isdir(currentFilePath):
			clear_despath(currentFilePath, endWithItems, exeludeItems)
		else:
			print "delete file "+currentFilePath
			os.remove(currentFilePath)
			
def copy_dirfile_to_dir(src_dir, des_dir, endWithItems, exeludeItems):
	print("call copy_dirfile_to_dir src_dir:"+src_dir+" des_dir:"+des_dir)
	try:
		os.makedirs(des_dir)
	except:
		pass
	for eachFile in os.listdir(src_dir):
		currentFilePath = os.path.join(src_dir, eachFile)
		currentDesFilePath = os.path.join(des_dir, eachFile)

		if os.path.isdir(currentFilePath):
			copy_dirfile_to_dir(currentFilePath, currentDesFilePath, endWithItems, exeludeItems)
		else:
			if endWithItems != None:
				if not is_endswith(eachFile, endWithItems):
					continue
			if exeludeItems != None:
				if is_endswith(eachFile, exeludeItems):
					continue
			print "copy file "+currentFilePath+" to "+ currentDesFilePath
#			delete_file(currentDesFilePath)
			shutil.copyfile(currentFilePath, currentDesFilePath)

def move_dirfile_to_dir(src_dir, des_dir, endWithItems, exeludeItems):		
	# try:
	# 	os.makedirs(des_dir)
	# except:
	# 	pass

	shutil.move(src_dir, des_dir)

if __name__ == "__main__":
	reload(sys)
	sys.setdefaultencoding("utf-8") 

	# result.put("path", path);
	# result.put("name", fileName);
				
#   通过编码的方法可以避免用这样的方式来写
	st1r = "你好"
	print type(st1r)
	print type(st1r.decode("utf-8"))
	print type(st1r.encode("gbk"))
	# .encode("gbk")
	# pass
#	print int(float(str(12)))
#	print int.parse("12.5")
#	util.recreate_dir(os.path.join(globalConfig.g_cache_dir, "compiled_lua"))
#	xmSrcPath = os.path.join(globalConfig.g_project_root_dir, "hall","xm")
#	outPutPath = os.path.join(globalConfig.g_cache_dir, "compiled_lua")
#	cmd = "%QUICK_V3_ROOT%//quick//bin//compile_scripts.bat -i "+ xmSrcPath + " -o " + outPutPath + " -m files -e xxtea_chunk -es ximi -ek HallArab20150301"
#	print cmd
#	os.system(cmd)
	# clear_despath(os.path.join(globalConfig.g_cache_dir, "zip_cache"), None, None)
	# shutil.rmtree(os.path.join(globalConfig.g_cache_dir, "zip_cache"))
	# move_dirfile_to_dir(os.path.join(globalConfig.g_cache_dir, "zip_cache"), "F:\opac", None, None)
	# print "%d"%(100)
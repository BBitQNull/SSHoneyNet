import os
import json
import stat

def is_special_dir(path):
    special_dirs = ['/proc', '/sys', '/dev', '/run']
    return any(path.startswith(d) for d in special_dirs)

def get_file_mode(path):
    try:
        st = os.lstat(path)
        if stat.S_ISDIR(st.st_mode):
            return "dir"
        elif stat.S_ISLNK(st.st_mode):
            return "link"
        elif stat.S_ISCHR(st.st_mode):
            return "char_device"
        elif stat.S_ISBLK(st.st_mode):
            return "block_device"
        elif stat.S_ISFIFO(st.st_mode):
            return "fifo"
        elif stat.S_ISSOCK(st.st_mode):
            return "socket"
        else:
            return "file"
    except Exception as e:
        return f"unknown ({str(e)})"

def process_special_dir(path):
    name = os.path.basename(path)
    result = {
        "name": name,
        "mode": "dynamic",
        "generator": name,
        "children": []
    }
    
    if path.startswith('/proc'):
        try:
            items = os.listdir(path)
            for item in items:
                if item.isdigit() or item in ['cpuinfo', 'meminfo', 'version']:
                    result["children"].append({
                        "name": item,
                        "mode": get_file_mode(os.path.join(path, item))
                    })
        except Exception as e:
            result["error"] = str(e)
    
    elif path.startswith('/sys'):
        try:
            items = os.listdir(path)
            for item in items[:10]:  # 限制数量
                result["children"].append({
                    "name": item,
                    "mode": get_file_mode(os.path.join(path, item))
                })
        except Exception as e:
            result["error"] = str(e)
    
    return result

def generate_file_tree(path, max_depth=3, current_depth=0):
    if current_depth > max_depth:
        return None
    
    name = os.path.basename(path) or '/'
    result = {"name": name}
    
    if is_special_dir(path):
        return process_special_dir(path)
    
    try:
        st = os.lstat(path)
        result["mode"] = get_file_mode(path)
        result["size"] = st.st_size
        result["mtime"] = st.st_mtime
        
        if result["mode"] == "link":
            result["target"] = os.readlink(path)
            return result
        
        if result["mode"] == "dir":
            children = []
            items = os.listdir(path)
            for item in items:
                item_path = os.path.join(path, item)
                child = generate_file_tree(item_path, max_depth, current_depth + 1)
                if child:
                    children.append(child)
            result["children"] = children
            
    except Exception as e:
        result["mode"] = "unknown"
        result["error"] = str(e)
    
    return result

# 生成完整文件树
root_tree = generate_file_tree("/", max_depth=4)  # 增加深度以获取更多细节

# 输出格式化的JSON
print(json.dumps(root_tree, indent=2))
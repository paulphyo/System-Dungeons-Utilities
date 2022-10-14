# Zennie's Receipt Helper 
# License © MIT
# Repository : https://github.com/z3nn13

import os
import json
from datetime import datetime
from colorama import Fore,Back,Style,init
from time import sleep

init(autoreset=True, convert=True)

# Function to clear screen
def cls():
    os.system('cls' if os.name=='nt' else 'clear')
cls()

# Function to return user id and directory
def get_user_input():
    data = {}
    menu = [
        "\t"+"*"*8+"MENU"+"*"*8 + "\n",
        "\t"+"-"*20
    ]
    print('\n'.join(menu))

    # Check for previous settings
    if os.path.exists('config.json'):
        with open('config.json', 'r', encoding="utf-8") as settings:
            data = json.load(settings)
            print('\n'.join(
                [
                    Fore.CYAN + "\tDetected Previous Settings",
                    Fore.CYAN + "\tUser ID | " + data['id'],
                    Fore.CYAN + "\tDirectory | " + data['directory'],
                    "",
                ]
            ))
            while True:
                reuse_settings = input(Fore.CYAN + "\tWould you like to reuse them? (Y/N):\n\t>").upper()
                if reuse_settings == 'Y':
                    return data
                elif reuse_settings == 'N':
                    print("\t"+"-"*20)
                    break

    # First time settings
    # Page 1
    data['id'] = input(Fore.GREEN + "\tEnter your Discord ID: (e.g. 408785106942164992)\n\t> ")
    cls()
    print(Style.RESET_ALL)
    # Page 2
    menu.pop(1) # removes a dash line below menu
    menu += [
        Fore.YELLOW + "\tSuccessfully stored ("+data['id']+") as user ID",
        Style.RESET_ALL + "\t"+"-"*20,
        "\tFor downloaded miner files,",
        "\t1. In different directory",
        "\t2. In current directory: {0}".format(os.getcwd()),
        ""
    ]
    print('\n'.join(menu))
    choice = int(input(Fore.GREEN + "\tEnter choice(1-2): "))
    data['directory'] = get_directory(choice)
    return data

def get_directory(choice):
    if choice == 2:
        return os.getcwd()
    elif choice == 1:
        while True:
            path = input(Fore.GREEN + "\tFull path to folder(example: C:/bachi/is/smexy): ")
            if os.path.isdir(path):
                return path
            print(Fore.RED + "\tPath does not exist. Please try again")

def save_settings_to_file(settings):
        with open('config.json', 'w', encoding="utf-8") as config:
            config.write(json.dumps(settings))
            config.close()
        
# Function to search and log crystals
def find_and_log_crystals(id,dirname):

    # Create folder if not exists
    outfolder_path = os.getcwd() + "\outfiles"
    os.makedirs(outfolder_path, exist_ok=True)

    # Create output file
    outfile_name = datetime.today().strftime('outfile_%b%d%Y.txt')
    outfile_path = os.path.join(outfolder_path, outfile_name)
    outfile = open(outfile_path, "w+", encoding='utf-8')
    
    crystals_found = []
    
    # Page 3
    cls()
    menu = [
        Back.LIGHTBLACK_EX + "*" + Back.CYAN + "User ID" + Style.RESET_ALL + " | " + Fore.YELLOW + id + Style.RESET_ALL,
        Back.LIGHTBLACK_EX + "*" + Back.CYAN + "Directory" + Style.RESET_ALL + " | " + Fore.YELLOW + dirname + Style.RESET_ALL,
        "-"*20, ""
        ]
    print("\n".join(menu))
    
    # Ordering Files By Date
    files = []
    for filename in os.listdir(dirname):
        if not filename.endswith('.txt'):
            continue
        if filename.startswith(('outfile', 'ReadMe')):
            continue
        files += [filename]
    try:
        sorted_files = sorted(files,key=lambda x: (datetime.strptime(x[:-4], "%b %d") or (datetime.strptime(x[:-4], "%B %d"))))
    except:
        sorted_files = files
        
    # Searching miner text files
    for file in sorted_files:
        with open(os.path.join(dirname,file), 'r', encoding="utf-8") as miner_json_file:
            users_array = json.load(miner_json_file)
            for user_dict in users_array:
                if user_dict['id'] == id:
                    print((Fore.GREEN + "Found" + Style.BRIGHT + " {0} " + Style.NORMAL + "with {1} crystals in '{2}'").format(user_dict['tag'], user_dict['quantity'], file))
                    outfile.write("%s : %s\n" % (file, json.dumps(user_dict, ensure_ascii=False)))
                    crystals_found += [user_dict['quantity']]
                    sleep(0.01)

    if not crystals_found:
        print(Fore.RED + "Searched %d files. 0 results found." % (len(files)))
    else:
        print(Fore.YELLOW + 'Writing to "' + outfile_path + '"...' )
        outText = "Total MC : %s = %s" % (" + ".join(str(crystals) for crystals in crystals_found), sum(crystals_found))
        outfile.write(outText)
        print(Fore.YELLOW + "All actions are done..!" )
        outfile.close()

    _ = str(input("Press enter to exit(): "))

def main():
    settings = get_user_input()
    save_settings_to_file(settings)
    find_and_log_crystals(settings['id'], settings['directory'])

main()
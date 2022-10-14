# Zennie's Receipt Helper
# License © MIT
# Repository : https://github.com/z3nn13

import os
import json
from datetime import datetime
from colorama import Fore, Back, Style, init
from time import sleep

# Refresh for colorama
init(autoreset=True, convert=True)

# Function to clear terminal
def cls():
    os.system("cls" if os.name == "nt" else "clear")


cls()


# Function to return user id and directory
def get_user_input():
    data = {}
    content = ["\t" + "*" * 8 + "MENU" + "*" * 8 + "\n", "\t" + "-" * 20]

    # Check for previous settings
    if not os.path.exists("config.json"):
        return show_page(1, content)
    else:
        return show_page(2, content)


# Function to return folder path
def get_directory(choice):
    if choice == 2:  # current_path
        return os.getcwd()
    elif choice == 1:  # different_path
        while True:
            path = input(
                Fore.GREEN + "\tFull path to folder(example: C:/bachi/is/smexy): "
            )
            if os.path.isdir(path):
                return path
            print(Fore.RED + "\tPath does not exist. Please try again")


# Function for display
def show_page(page, content):
    if page == 1:
        print("\n".join(content))
        data["id"] = input(
            Fore.GREEN + "\tEnter a Discord ID: (e.g. 408785106942164992)\n\t> "
        )
        print(Style.RESET_ALL)
        cls()
        content.pop(1)  # removes a dash line below menu
        content += [
            Fore.YELLOW + "\tSuccessfully stored (" + data["id"] + ") as user ID",
            Fore.WHITE + "\t" + "-" * 20,
            "\tFor downloaded miner files,",
            "\t1. In different directory",
            "\t2. In current directory: {0}".format(os.getcwd()),
            "",
        ]
        print("\n".join(content))
        choice = int(input(Fore.GREEN + "\tEnter choice(1-2): "))
        data["directory"] = get_directory(choice)
        return data

    elif page == 2:
        print("\n".join(content))
        with open("config.json", "r", encoding="utf-8") as settings:
            data = json.load(settings)
            print(
                "\n".join(
                    [
                        Fore.CYAN + "\tDetected Previous Settings",
                        Fore.CYAN + "\tUser ID | " + data["id"],
                        Fore.CYAN + "\tDirectory | " + data["directory"],
                        "",
                    ]
                )
            )
            while True:
                reuse_settings = input(
                    Fore.CYAN + "\tWould you like to reuse them? (Y/N):\n\t>"
                ).upper()
                if reuse_settings == "Y":
                    return data
                elif reuse_settings == "N":
                    print("\t" + "-" * 20)
                    show_page(1)
    elif page == 3:
        cls()
        header = [
            # f"{Back.LIGHTBLACK_EX}*{Back.CYAN} User ID {Style.RESET_ALL}|{Fore.YELLOW} {content['id']} {Style.RESET_ALL}",
            f"*User ID | {content['id']}",
            f"*Directory | {content['dirname']}",
            "-" * 20,
            "",
        ]
        print("\n".join(header))


# Function to search and log crystals
def find_and_log_crystals(settings):
    id = settings["id"]
    dirname = settings["directory"]

    # Create folder if not exists
    outfolder_path = os.getcwd() + "\outfiles"
    os.makedirs(outfolder_path, exist_ok=True)

    # Create output file
    outfile_name = datetime.today().strftime("outfile_%b%d%Y.txt")
    outfile_path = os.path.join(outfolder_path, outfile_name)

    crystals_found = []

    # Page 3
    show_page(3, settings)

    # Ordering Files By Date
    files = []
    for file in os.listdir(dirname):
        if not file.endswith(".txt"):
            continue
        if file.startswith(("outfile", "ReadMe")):
            continue
        files += [file]
    try:
        sorted_files = sorted(files, key=lambda x: (datetime.strptime(x[:-4], "%b %d")))
    except ValueError:
        sorted_files = files

    # Searching miner text files
    with open(outfile_path, "w+", encoding="utf-8") as outfile:
        for filename in sorted_files:
            with open(
                os.path.join(dirname, filename), "r", encoding="utf-8"
            ) as miner_json_file:
                users_array = json.load(miner_json_file)
                for user_dict in users_array:
                    if user_dict["id"] == id:
                        print(
                            f"%s Found %s {user_dict['tag']} %s with {user_dict['quantity']} crystals in '{filename}'"
                            % [Fore.GREEN, Style.BRIGHT, Style.NORMAL]
                        )
                        outfile.write(
                            f"{filename} : {json.dumps(user_dict, ensure_ascii=False)}"
                        )
                        crystals_found += [user_dict["quantity"]]
                        sleep(0.01)

        if not crystals_found:
            print(
                "{} Searched {filecount} files. 0 results found.".format(
                    Fore.RED, len(files)
                )
            )
        else:
            print(f"{Fore.YELLOW} Writing to {outfile_path} ...")
            outfile.write(
                "Total MC : {additions} = {sum} ".format(
                    " + ".join(str(crystals) for crystals in crystals_found),
                    sum(crystals_found),
                )
            )
            print(f'{Fore.YELLOW} + "All actions are done..!"')
            outfile.close()

    _ = str(input("Press enter to exit(): "))


def save_settings_to_file(settings):
    with open("config.json", "w", encoding="utf-8") as config:
        config.write(json.dumps(settings))
        config.close()


def main():
    settings = get_user_input()
    save_settings_to_file(settings)
    find_and_log_crystals(settings)


if __name__ == "__main":
    main()

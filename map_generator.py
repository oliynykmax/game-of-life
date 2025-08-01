import random
import os

def generate_3_gol_files(n, n2, fill_percentage, directory="maps"):
    os.makedirs(directory, exist_ok=True)

    for idx in range(3):
        filename = f"{n}x{n2}_{idx}_{fill_percentage}%"
        filepath = os.path.join(directory, filename)

        with open(filepath, 'w') as f:
            for _ in range(n2):
                row = ''.join(random.choices(['.', 'X'], weights=[100 - fill_percentage, fill_percentage], k=n))
                f.write(row + '\n')

    print(f"Generated 3 files in folder '{directory}' with {n}x{n2} grids, each having approximately {fill_percentage}% live cells.")

if __name__ == "__main__":
    n = int(input("Enter the value of n: "))
    n2 = int(input("Enter the value of n2: "))
    fill_percentage = float(input("Enter the fill percentage (0 to 100): "))
    generate_3_gol_files(n, n2, fill_percentage)

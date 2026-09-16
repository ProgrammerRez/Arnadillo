import pandas as pd


df = pd.read_csv("Data/customer_master.csv")
print(df.describe())
print("==="*30)
print(df.info())
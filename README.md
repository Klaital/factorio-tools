# factorio-tools
Tools for analyzing Factorio game data

# Requirements

1. Ingest a JSON database of the items in game and the recipes required to construct them. It should validate that all recipes reference real items.
2. Given an Item name, generate a list of the quantities of each component item that will be needed to craft it. This is to calculate requirements like "I want to make a Steam Engine. What all do I need to gather/craft first?"
3. Given an Item and Rate, generate a list of the quantities of each component item that will need to be crafted per minute to meet the Rate. This is to calculate requirements like, "I want to craft 1k processing units per hour. What do I need?"


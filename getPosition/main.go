package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type Tick struct {
	Mag  *felt.Felt
	Sign *felt.Felt
}

type MintEventStruct struct {
	Sender    *felt.Felt
	Owner     *felt.Felt
	TickLower Tick
	TickUpper Tick
	Amount    *felt.Felt
	Amount0   *felt.Felt
	Amount1   *felt.Felt
}
type CombinedStruct struct {
	MintEvent MintEventStruct `json:"mintEvent"`
	Position  PositionInfo    `json:"position"`
}

type PositionInfo struct {
	Liquidity                *felt.Felt
	FeeGrowthInside0LastX128 *felt.Felt
	FeeGrowthInside1LastX128 *felt.Felt
	TokensOwed0              *felt.Felt
	TokensOwed1              *felt.Felt
}

func QueryIncreaseLiquidityEvents(startblock uint64, endblock uint64) {
	base := "https://starknet-mainnet.public.blastapi.io"
	fmt.Println("Starting simpeCall example")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, _ := utils.HexToFelt("0x00287d2ff1c39a44cd18d9dc7ed5617c9cb16b65090db6a0f689aa14755e4e5e")
	IncreaseLiquidity, _ := utils.HexToFelt("0x3159b5cf425448a640d4666273b957d14e6c7a6a5a6a758af3c58d7ab7841fb")

	eventInput := rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.BlockID{
				Number: &startblock,
			},
			ToBlock: rpc.BlockID{
				Number: &endblock,
			},
			Address: contractAddress,
			Keys: [][]*felt.Felt{{
				IncreaseLiquidity,
			}},
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	}

	events, err := clientv02.Events(context.Background(), eventInput)
	if err != nil {
		fmt.Println("Unsuccessful")
		panic(err)
	}

	for _, emittedEvent := range events.Events {
		fmt.Println("\n", emittedEvent.BlockNumber, " Block Number")
		token_id := emittedEvent.Event.Data[0].String()
		fmt.Println(token_id, " Token ID")
		liquidity := emittedEvent.Event.Data[2].String()
		fmt.Println(liquidity, " Liquidity")
		amount0 := emittedEvent.Event.Data[3].String()
		fmt.Println(amount0, " Liquidity0")
		amount1 := emittedEvent.Event.Data[5].String()
		fmt.Println(amount1, " Liquidity1")

	}
}

func QueryMintEvents(startblock uint64, endblock uint64) []MintEventStruct {
	base := "https://starknet-mainnet.public.blastapi.io"
	fmt.Println("Starting simpeCall example")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client, did you specify the url in the .env.mainnet?")
		panic(err)
	}
	clientv02 := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, _ := utils.HexToFelt("0x06096f2a295571bd45b627c88c7f1760cc8ff27c1c3c204e68ed3fe040844d2d")

	Mint, _ := utils.HexToFelt("0x34e55c1cd55f1338241b50d352f0e91c7e4ffad0e4271d64eb347589ebdfd16")

	eventInput := rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: rpc.BlockID{
				Number: &startblock,
			},
			ToBlock: rpc.BlockID{
				Number: &endblock,
			},
			Address: contractAddress,
			Keys: [][]*felt.Felt{{
				Mint,
			}},
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ChunkSize: 1000,
		},
	}

	var mintEvents []MintEventStruct

	events, err := clientv02.Events(context.Background(), eventInput)
	if err != nil {
		fmt.Println("Unsuccessful")
		panic(err)
	}

	for _, emittedEvent := range events.Events {
		var event MintEventStruct

		event.Sender = emittedEvent.Event.Data[0]
		event.Owner = emittedEvent.Event.Data[1]
		event.TickLower.Mag = emittedEvent.Event.Data[2]
		event.TickLower.Sign = emittedEvent.Event.Data[3]
		event.TickUpper.Mag = emittedEvent.Event.Data[4]
		event.TickUpper.Sign = emittedEvent.Event.Data[5]
		event.Amount = emittedEvent.Event.Data[6]
		event.Amount0 = emittedEvent.Event.Data[7]
		event.Amount1 = emittedEvent.Event.Data[9]

		mintEvents = append(mintEvents, event)
	}

	return mintEvents
}

func get_position_info(owners, tick_lower_mags, tick_lower_signs, tick_upper_mags, tick_upper_signs []*felt.Felt) []PositionInfo {
	base := "https://starknet-mainnet.public.blastapi.io"
	fmt.Println("Starting batch call example")
	c, err := ethrpc.DialContext(context.Background(), base)
	if err != nil {
		fmt.Println("Failed to connect to the client")
		panic(err)
	}
	client := rpc.NewProvider(c)
	fmt.Println("Established connection with the client")

	contractAddress, _ := utils.HexToFelt("0x06096f2a295571bd45b627c88c7f1760cc8ff27c1c3c204e68ed3fe040844d2d")

	var positions []PositionInfo

	if len(owners) == len(tick_lower_mags) && len(owners) == len(tick_lower_signs) && len(owners) == len(tick_upper_mags) && len(owners) == len(tick_upper_signs) {
		for i := range owners {
			var position PositionInfo
			args := []*felt.Felt{owners[i], tick_lower_mags[i], tick_lower_signs[i], tick_upper_mags[i], tick_upper_signs[i]}

			tx := rpc.FunctionCall{
				ContractAddress:    contractAddress,
				EntryPointSelector: utils.GetSelectorFromNameFelt("get_position_info"), // Ensure correct selector usage
				Calldata:           args,
			}

			callResponse, err := client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
			if err != nil {
				fmt.Println("Error calling smart contract:", err)
				panic(err)
			}

			position.Liquidity = callResponse[0]
			position.FeeGrowthInside0LastX128 = callResponse[1]
			position.FeeGrowthInside1LastX128 = callResponse[3]
			position.TokensOwed0 = callResponse[5]
			position.TokensOwed1 = callResponse[6]

			positions = append(positions, position)
		}
	} else {
		fmt.Println("Error: Input slices do not have the same length.")
	}

	return positions
}

func WriteDataToJsonFile(data []CombinedStruct, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var startblock_router, endblock_router uint64

	startblock_router = 550860
	endblock_router = 550862

	println("QueryMintEvents Events")
	res := QueryMintEvents(startblock_router, endblock_router)

	var owners, tickLowerMags, tickLowerSigns, tickUpperMags, tickUpperSigns []*felt.Felt
	for _, event := range res {
		owners = append(owners, event.Owner)
		tickLowerMags = append(tickLowerMags, event.TickLower.Mag)
		tickLowerSigns = append(tickLowerSigns, event.TickLower.Sign)
		tickUpperMags = append(tickUpperMags, event.TickUpper.Mag)
		tickUpperSigns = append(tickUpperSigns, event.TickUpper.Sign)
	}

	positionInfos := get_position_info(owners, tickLowerMags, tickLowerSigns, tickUpperMags, tickUpperSigns)

	var combinedData []CombinedStruct

	if len(res) == len(positionInfos) {
		for i, mintEvent := range res {
			combined := CombinedStruct{
				MintEvent: mintEvent,
				Position:  positionInfos[i],
			}
			combinedData = append(combinedData, combined)
		}
	} else {
		fmt.Println("The length of mint events and position info slices do not match.")
		return
	}

	filenameCombined := "combined_output.json"
	err := WriteDataToJsonFile(combinedData, filenameCombined)
	if err != nil {
		fmt.Println("Failed to write combined data to JSON file:", err)
		return
	}

	fmt.Println("Combined data successfully written to", filenameCombined)

}

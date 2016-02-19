package main



func CountIslands(input [][]int) int {
	if len(input)==0 || len(input[0])==0{
     return 0;
    }
    count:= 0;
 
    for i:=0; i<len(input); i++ {
        for j:=0; j<len(input[0]); j++ {
            if input[i][j] == 1 {
                count++;
                merge(input, i, j);
            }
        }
    }
    return count;
}


func merge(input [][]int, i int, j int){
    
    if i<0 || j<0 || i>len(input)-1 || j>len(input[0])-1 {
        return ;
    }
    //if current cell is water or visited
    if input[i][j] != 1 {
        return;
 }
    //set visited cell to '0'
    input[i][j] =0;

     
    //merge all adjacent land
    merge(input, i-1, j);
    merge(input, i+1, j);
    merge(input, i, j-1);
    merge(input, i, j+1);
}
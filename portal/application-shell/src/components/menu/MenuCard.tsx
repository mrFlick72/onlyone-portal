import React from 'react';
import {Card, CardActionArea, CardContent, CardHeader, Grid2} from "@mui/material";

const classes = {
    root: {
        maxWidth: 345
    },
    media: {
        height: 140
    }
}

interface MenuCardProps {
    title?: string
    content: any
    linkTo: string
}

const MenuCard: React.FC<MenuCardProps> = ({title, content, linkTo}) => {
    return (
        <Grid2 size={{ xs: 12, sm: 6, md: 4, lg: 4, xl: 4 }}>
            <a href={linkTo} style={{textDecoration: 'none', textAlign: 'center'}}>
                <Card style={classes.root}>
                    <CardActionArea>
                        <CardHeader
                            style={classes.media}
                            title={title}
                        />
                        <CardContent>
                            {content}
                        </CardContent>
                    </CardActionArea>
                </Card>
            </a>
        </Grid2>
    );
}

export default MenuCard